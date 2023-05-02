package service

import (
	"context"
	"fmt"
	"time"

	rpc "github.com/sirjager/rpcs/trueauth/go"
	"github.com/sirjager/trueauth/db/sqlc"
	"github.com/sirjager/trueauth/mail"
	"github.com/sirjager/trueauth/tokens"
	"github.com/sirjager/trueauth/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Request email verification and verify email
//
// 1. if CODE is provided then verify code
// 2. If NO CODE then send email verification code
func (s *TrueAuthService) Delete(ctx context.Context, req *rpc.DeleteRequest) (*rpc.DeleteResponse, error) {
	authorized, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthenticatedError(err)
	}

	// If email is not verified then we will simply delete the account
	if !authorized.User.EmailVerified {
		if err := s.store.Delete_User(ctx, authorized.User.ID); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to terminate : %s", err.Error())
		}
		return &rpc.DeleteResponse{Message: "you account has been terminated"}, nil
	}

	meta := s.extractMetadata(ctx)

	//? if no code is provided means, account is requesting deletion code
	if req.GetCode() == "" {
		// If code is sent recently then wait for code request cooldown
		if time.Since(authorized.User.LastDeleteSentAt) < s.config.DeleteTokenCooldown {
			tryAfter := time.Until(authorized.User.LastDeleteSentAt.Add(s.config.DeleteTokenCooldown))
			return &rpc.DeleteResponse{
				Message: fmt.Sprintf("account deletion code has been requested recently, please try again after %s", tryAfter),
			}, nil
		}

		// generate a random code
		deletionCode := utils.RandomNumberAsString(6)
		durationTTL := s.config.DeleteTokenTTL

		saveDeletionToken, _, err := s.tokens.CreateToken(
			tokens.PayloadData{
				UserID:       authorized.User.ID,
				UserEmail:    authorized.User.Email,
				DeletionCode: deletionCode,
			}, durationTTL,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create code: %s", err.Error())
		}

		mail := mail.Mail{To: []string{authorized.User.Email}}
		mail.Subject = "Thank you for joining us. Please confirm your mail"
		mail.Body = fmt.Sprintf(`
		Hello <br/>
		Account deletion requested from :
		IP Address 	   :  <b>%s</b> <br/>
		User Agent 	   :  <b>%s</b> <br/>
		Deletion Code  :  <b>%s</b> <br/>
		Valid Only For :  <b>%s</b> <br/>
		If you didn't request this, simply ignore this message. <br/> <br/>
		Thank You`,
			meta.ClientIp, meta.UserAgent, deletionCode, durationTTL.String(),
		)

		if err = s.store.Update_User_DeleteTokenTx(ctx, sqlc.Update_User_DeleteTokenTxParams{
			Update_User_DeleteTokenParams: sqlc.Update_User_DeleteTokenParams{
				ID:               authorized.User.ID,
				DeleteToken:      saveDeletionToken,
				LastDeleteSentAt: time.Now(),
			},
			BeforeUpdate: func() error {
				return s.mailer.SendMail(mail)
			},
		}); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to initiate deletion processs: %s", err.Error())
		}

		return &rpc.DeleteResponse{Message: "deletion code has been sent to your email"}, nil
	}

	//? if code is provided means account is submiting code
	violations := validateDeleteRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}

	deleteTokenPayload, err := s.tokens.VerifyToken(authorized.User.DeleteToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid email verification code: %s", err.Error())
	}

	if deleteTokenPayload.Payload.UserEmail != authorized.User.Email {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched email")
	}

	if deleteTokenPayload.Payload.UserID.String() != authorized.User.ID.String() {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched account id")
	}

	if deleteTokenPayload.Payload.DeletionCode != req.GetCode() {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched code")
	}

	// Now delete the account
	if err := s.store.Delete_User(ctx, authorized.User.ID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to terminate : %s", err.Error())
	}
	return &rpc.DeleteResponse{Message: "you account has been terminated"}, nil
}

func validateDeleteRequest(req *rpc.DeleteRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if len(req.GetCode()) != 6 {
		violations = append(violations, fieldViolation("code", fmt.Errorf("invalid code")))
	}
	return
}
