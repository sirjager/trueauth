package service

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rpc "github.com/sirjager/rpcs/trueauth/go"

	"github.com/sirjager/trueauth/internal/db/sqlc"
	"github.com/sirjager/trueauth/internal/worker"

	"github.com/sirjager/trueauth/pkg/mail"
	"github.com/sirjager/trueauth/pkg/tokens"
	"github.com/sirjager/trueauth/pkg/utils"
)

// Request email verification and verify email
//
// 1. if CODE is provided then verify code
// 2. If NO CODE then send email verification code
func (s *CoreService) Verify(ctx context.Context, req *rpc.VerifyRequest) (*rpc.VerifyResponse, error) {
	authorized, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthenticatedError(err)
	}

	if authorized.User.EmailVerified {
		return &rpc.VerifyResponse{Message: fmt.Sprintf("email %s is already verified", authorized.User.Email)}, nil
	}

	//? if no code is provided means, user is requesting email verification code
	if req.GetCode() == "" {

		// If verification code is sent recently then wait for verification request cooldown
		if time.Since(authorized.User.LastVerifySentAt) < s.Config.VerifyTokenCooldown {
			tryAfter := time.Until(authorized.User.LastVerifySentAt.Add(s.Config.VerifyTokenCooldown))
			return &rpc.VerifyResponse{
				Message: fmt.Sprintf("email verification has been requested recently, please try again after %s", tryAfter),
			}, nil
		}

		// generate a random code
		verificationCode := utils.RandomNumberAsString(6)
		durationTTL := s.Config.VerifyTokenTTL

		verificationToken, _, err := s.tokens.CreateToken(
			tokens.PayloadData{
				UserID:           authorized.User.ID,
				UserEmail:        authorized.User.Email,
				VerificationCode: verificationCode,
			}, durationTTL,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create code: %s", err.Error())
		}

		mail := mail.Mail{To: []string{authorized.User.Email}}
		mail.Subject = "Email verification code"
		mail.Body = fmt.Sprintf(`
		Hello <br/>
		Your email verification code is : <b>%s</b> <br/>
		This code is only valid for %s <br/> <br/>
		If you didn't request this, simply ignore this message. <br/> <br/>
		Thank You`, verificationCode, durationTTL.String())

		if err = s.store.UpdateUserVerifyTokenTx(ctx, sqlc.UpdateUserVerifyTokenTxParams{
			UpdateUserVerifyTokenParams: sqlc.UpdateUserVerifyTokenParams{
				LastVerifySentAt: time.Now(),
				VerifyToken:      verificationToken,
				ID:               authorized.User.ID,
			},
			BeforeUpdate: func() error {
				return s.mailer.SendMail(mail)
			},
		}); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update email confirmation token: %s", err.Error())
		}

		return &rpc.VerifyResponse{Message: fmt.Sprintf("email verification code sent to your email %s", authorized.User.Email)}, nil
	}

	//? if code is provided means user is submiting email verification code

	violations := validateVerifyRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}

	tokenPayload, err := s.tokens.VerifyToken(authorized.User.VerifyToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid email verification code: %s", err.Error())
	}

	if tokenPayload.Payload.UserEmail != authorized.User.Email {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched email")
	}

	if tokenPayload.Payload.UserID.String() != authorized.User.ID.String() {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched user id")
	}

	if tokenPayload.Payload.VerificationCode != req.GetCode() {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched code")
	}

	if err = s.store.UpdateUserEmailVerifiedTx(ctx, sqlc.UpdateUserEmailVerifiedTxParams{
		UpdateUserEmailVerifiedParams: sqlc.UpdateUserEmailVerifiedParams{
			EmailVerified: true,
			VerifyToken:   "null",
			ID:            authorized.User.ID,
		},
		AfterUpdate: func() error {
			opts := []asynq.Option{
				asynq.MaxRetry(5),
				asynq.Group(worker.QUEUE_LOW),
				asynq.ProcessIn(time.Second * 10),
			}
			return s.taskDistributor.DistributeTaskSendEmailVerified(ctx,
				worker.PayloadSendEmailVerified{Email: authorized.User.Email}, opts...)
		},
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set email verified: %s", err.Error())
	}

	return &rpc.VerifyResponse{Message: fmt.Sprintf("email %s successfully verified", authorized.User.Email)}, nil
}

func validateVerifyRequest(req *rpc.VerifyRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if len(req.GetCode()) != 6 {
		violations = append(violations, fieldViolation("code", fmt.Errorf("invalid code")))
	}
	return
}
