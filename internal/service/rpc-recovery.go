package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rpc "github.com/sirjager/rpcs/trueauth/go"

	"github.com/sirjager/trueauth/internal/db/sqlc"

	"github.com/sirjager/trueauth/pkg/crypt"
	"github.com/sirjager/trueauth/pkg/mail"
	"github.com/sirjager/trueauth/pkg/tokens"
	"github.com/sirjager/trueauth/pkg/utils"
	"github.com/sirjager/trueauth/pkg/validator"
)

// Request email verification and verify email
//
// 1. if CODE is provided then verify code
// 2. If NO CODE then send email verification code
func (s *CoreService) Recovery(ctx context.Context, req *rpc.RecoveryRequest) (*rpc.RecoveryResponse, error) {
	violations := validateRecoveryRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}

	user, err := s.store.ReadUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "no such email %s ", req.GetEmail())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if s.isUnKnownIP(ctx, user) {
		return nil, unknownIPError()
	}

	meta := s.extractMetadata(ctx)

	if req.Code == "" {
		// If password recovery code is sent recently then wait for recovery code request cooldown
		if time.Since(user.LastRecoverySentAt) < s.Config.ResetTokenCooldown {
			tryAfter := time.Until(user.LastRecoverySentAt.Add(s.Config.ResetTokenCooldown))
			return &rpc.RecoveryResponse{
				Message: fmt.Sprintf("password recovery code has been requested recently, please try again after %s", tryAfter),
			}, nil
		}

		// generate a random code
		recoveryCode := utils.RandomNumberAsString(6)
		durationTTL := s.Config.VerifyTokenTTL

		saveRecoveryToken, _, err := s.tokens.CreateToken(tokens.PayloadData{
			UserID:       user.ID,
			UserEmail:    user.Email,
			RecoveryCode: recoveryCode,
		}, durationTTL)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create code: %s", err.Error())
		}

		mail := mail.Mail{To: []string{user.Email}}
		mail.Subject = "Thank you for joining us. Please confirm your mail"
		mail.Body = fmt.Sprintf(`
		Hello <br/>
		<b> Someone has requested password reset code from </b> <br/>
		IP Address : <b>%s</b> <br/>
		User Agent : <b>%s</b> <br/>
		If this was you, Use the password reset code : <b> %s </b> to change password. <br/>
		If this wans't you then you can simply ignore the request. <br/>
		This code is only valid for %s <br/>
		Thank You`,
			meta.ClientIp, meta.UserAgent, recoveryCode, durationTTL.String(),
		)

		if err = s.store.UpdateUserRecoveryTokenTx(ctx, sqlc.UpdateUserRecoveryTokenTxParams{
			UpdateUserRecoveryTokenParams: sqlc.UpdateUserRecoveryTokenParams{
				LastRecoverySentAt: time.Now(),
				ID:                 user.ID,
				RecoveryToken:      saveRecoveryToken,
			},
			BeforeUpdate: func() error {
				return s.mailer.SendMail(mail)
			},
		}); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to fullfill request for password recovery : %s", err.Error())
		}

		return &rpc.RecoveryResponse{
			Message: fmt.Sprintf("password recovery code has been sent to your email %s", user.Email),
		}, nil
	}

	// If code is available means we just need to verify code and generate new password
	recoveryPayload, err := s.tokens.VerifyToken(user.RecoveryToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid email verification code: %s", err.Error())
	}

	if recoveryPayload.Payload.UserEmail != user.Email {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched email")
	}

	if recoveryPayload.Payload.UserID.String() != user.ID.String() {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched user id")
	}

	if recoveryPayload.Payload.RecoveryCode != req.GetCode() {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched code")
	}

	// Now generate a new password, update in database and send it to user's email

	newRandomPassword := utils.RandomPassword()
	hasedRandomPassword, err := crypt.HashPassword(newRandomPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err.Error())
	}

	// We will send new password to user's email
	mail := mail.Mail{To: []string{user.Email}}
	mail.Subject = "Thank you for joining us. Please confirm your mail"
	mail.Body = fmt.Sprintf(`
		Hello <br/>
		You password has been reset. </br>
		IP Address 	 : <b>%s</b> <br/>
		User Agent 	 : <b>%s</b> <br/>
		New Password : <b>%s</b> <br/>
		Thank You`,
		meta.ClientIp, meta.UserAgent, newRandomPassword,
	)

	if err = s.store.UpdateUserResetPasswordTx(ctx, sqlc.UpdateUserResetPasswordTxParams{
		UpdateUserResetPasswordParams: sqlc.UpdateUserResetPasswordParams{
			RecoveryToken: "null",
			ID:            user.ID,
			Password:      hasedRandomPassword,
		},
		BeforeUpdate: func() error {
			return s.mailer.SendMail(mail)
		},
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reset password: %s", err.Error())
	}

	return &rpc.RecoveryResponse{Message: "your new password has been sent to your email"}, nil
}

func validateRecoveryRequest(req *rpc.RecoveryRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	return
}
