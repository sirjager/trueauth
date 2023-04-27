package service

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	rpc "github.com/sirjager/rpcs/trueauth/go"
	"github.com/sirjager/trueauth/db/sqlc"
	"github.com/sirjager/trueauth/mail"
	"github.com/sirjager/trueauth/tokens"
	"github.com/sirjager/trueauth/utils"
	"github.com/sirjager/trueauth/worker"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Request email verification and verify email
//
// 1. if CODE is provided then verify code
// 2. If NO CODE then send email verification code
func (s *TrueAuthService) VerifyEmail(ctx context.Context, req *rpc.VerifyEmailRequest) (*rpc.VerifyEmailResponse, error) {
	authorized, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthenticatedError(err)
	}

	if authorized.Account.EmailVerified {
		return &rpc.VerifyEmailResponse{Message: fmt.Sprintf("email %s is already verified", authorized.Account.Email)}, nil
	}

	//? if no code is provided means, account is requesting email verification code
	if req.GetCode() == "" {

		// If verification code is sent recently then wait for verification request cooldown
		if time.Since(authorized.Account.LastConfirmationSentAt) < s.config.VerifyTokenCooldown {
			tryAfter := time.Until(authorized.Account.LastConfirmationSentAt.Add(s.config.VerifyTokenCooldown))
			return &rpc.VerifyEmailResponse{
				Message: fmt.Sprintf("email verification has been requested recently, please try again after %s", tryAfter),
			}, nil
		}

		// generate a random code
		sixDigitCode := utils.RandomNumberAsString(6)
		durationTTL := s.config.VerifyTokenTTL

		token, _, err := s.tokens.CreateToken(tokens.PayloadData{
			AccountID:             authorized.Account.ID,
			AccountEmail:          authorized.Account.Email,
			EmailVerificationCode: sixDigitCode,
		}, durationTTL)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create code: %s", err.Error())
		}

		mail := mail.Mail{To: []string{authorized.Account.Email}}
		mail.Subject = "Thank you for joining us. Please confirm your mail"
		mail.Body = fmt.Sprintf(`
		Hello <br/>
		Your email verification code is : <b>%s</b> <br/>
		This code is only valid for %s <br/> <br/>
		If you didn't request this, simply ignore this message. <br/> <br/>
		Thank You`, sixDigitCode, durationTTL.String())

		if err = s.store.UpdateAccountEmailConfirmationTokenTx(ctx, sqlc.UpdateAccountEmailConfirmationTokenTxParams{
			UpdateAccountEmailConfirmationTokenParams: sqlc.UpdateAccountEmailConfirmationTokenParams{
				ID:                     authorized.Account.ID,
				ConfirmationToken:      token,
				LastConfirmationSentAt: time.Now(),
			},
			BeforeUpdate: func() error {
				return s.mailer.SendMail(mail)
			},
		}); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update email confirmation token: %s", err.Error())
		}

		return &rpc.VerifyEmailResponse{Message: fmt.Sprintf("email verification code sent to your email %s", authorized.Account.Email)}, nil
	}

	//? if code is provided means account is submiting email verification code

	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}

	tokenPayload, err := s.tokens.VerifyToken(authorized.Account.ConfirmationToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid email verification code: %s", err.Error())
	}

	if tokenPayload.Payload.AccountEmail != authorized.Account.Email {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched email")
	}

	if tokenPayload.Payload.AccountID.String() != authorized.Account.ID.String() {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched account id")
	}

	if tokenPayload.Payload.EmailVerificationCode != req.GetCode() {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched code")
	}

	if err = s.store.UpdateAccountEmailVerifiedTx(ctx, sqlc.UpdateAccountEmailVerifiedTxParams{
		UpdateAccountEmailVerifiedParams: sqlc.UpdateAccountEmailVerifiedParams{
			EmailVerified:     true,
			ConfirmationToken: "null",
			ID:                authorized.Account.ID,
		},
		AfterUpdate: func() error {
			opts := []asynq.Option{
				asynq.MaxRetry(5),
				asynq.Group(worker.QUEUE_LOW),
				asynq.ProcessIn(time.Second * 10),
			}
			return s.taskDistributor.DistributeTaskSendEmailVerified(ctx,
				worker.PayloadSendEmailVerified{Email: authorized.Account.Email}, opts...)
		},
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set email verified: %s", err.Error())
	}

	return &rpc.VerifyEmailResponse{Message: fmt.Sprintf("email %s successfully verified", authorized.Account.Email)}, nil
}

func validateVerifyEmailRequest(req *rpc.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if len(req.GetCode()) != 6 {
		violations = append(violations, fieldViolation("code", fmt.Errorf("invalid code")))
	}
	return
}
