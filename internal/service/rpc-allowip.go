package service

import (
	"context"
	"fmt"
	"time"

	rpc "github.com/sirjager/rpcs/trueauth/go"
	"github.com/sirjager/trueauth/internal/db/sqlc"
	"github.com/sirjager/trueauth/pkg/mail"
	"github.com/sirjager/trueauth/pkg/tokens"
	"github.com/sirjager/trueauth/pkg/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreService) AllowIP(ctx context.Context, req *rpc.AllowIPRequest) (*rpc.AllowIPResponse, error) {
	// 1. we will call authorization, if everything is fine we will get authorized user.
	authorized, _ := s.authorize(ctx)

	meta := s.extractMetadata(ctx)

	// if there is no code means user is trying to whitelist new ip address
	if req.GetCode() == "" {
		// In order to send send emails we will see to it that email is verified.
		if !authorized.User.EmailVerified {
			return &rpc.AllowIPResponse{Message: "you need to verify your email first"}, nil
		}

		// we will also check and make sure email and user id is not empty,
		// since we are not checking errors from step 1. authorized function
		// we did not check error becoz authorized func will always return error since we are also checking for unknown ips.
		if authorized.User.ID.String() == "" || authorized.User.Email == "" {
			return nil, unAuthenticatedError(fmt.Errorf("invalid email or user id"))
		}

		if !s.isUnKnownIP(ctx, authorized.User) {
			return &rpc.AllowIPResponse{Message: "ip address is already in whitelist"}, nil
		}

		// If allow ip code is sent recently then wait for allow ip code request cooldown
		if time.Since(authorized.User.LastAllowipSentAt) < s.Config.AllowIPTokenCooldown {
			tryAfter := time.Until(authorized.User.LastAllowipSentAt.Add(s.Config.AllowIPTokenCooldown))
			return &rpc.AllowIPResponse{
				Message: fmt.Sprintf("code to allow your current ip address has already been sent recently, please try again after %s", tryAfter),
			}, nil
		}

		// generate a random code
		allowIPCode := utils.RandomNumberAsString(6)
		durationTTL := s.Config.AllowIPTokenTTL

		saveAllowIPToken, _, err := s.tokens.CreateToken(tokens.PayloadData{
			UserID:      authorized.User.ID,
			UserEmail:   authorized.User.Email,
			AllowIPCode: allowIPCode,
			AllowIP:     meta.ClientIp,
		}, durationTTL)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create code: %s", err.Error())
		}

		mail := mail.Mail{To: []string{authorized.User.Email}}
		mail.Subject = "Request from unknown ip address from your account"
		mail.Body = fmt.Sprintf(`
		Hello <br/>
		<b> Someone has requested to whitelist a unknown ip address </b> <br/>
		Unknown IP Address : <b>%s</b> <br/>
		User Agent : <b>%s</b> <br/>
		If you wish to allow requests from this ip address, 
		Use the allow ip adress code : <b> %s </b> to whitelist this ip address. <br/>
		If this wans't you then you can simply ignore the request. <br/>
		This code is only valid for %s <br/>
		Thank You`,
			meta.ClientIp, meta.UserAgent, allowIPCode, durationTTL.String(),
		)

		if err = s.store.UpdateUserAllowIPTokenTx(ctx, sqlc.UpdateUserAllowIPTokenTxParams{
			UpdateUserAllowIPTokenParams: sqlc.UpdateUserAllowIPTokenParams{
				ID:                authorized.User.ID,
				LastAllowipSentAt: time.Now(),
				AllowipToken:      saveAllowIPToken,
			},
			BeforeUpdate: func() error {
				return s.mailer.SendMail(mail)
			},
		}); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to fullfill request for allowip new ip adress : %s", err.Error())
		}

		return &rpc.AllowIPResponse{
			Message: "code has been sent to your email to allow requests from this ip address",
		}, nil
	}

	// from here we will always have allow ip code.
	violations := validateAllowIPRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}

	tokenPayload, err := s.tokens.VerifyToken(authorized.User.AllowipToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid code: %s", err.Error())
	}

	if tokenPayload.Payload.UserEmail != authorized.User.Email {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched email")
	}

	if tokenPayload.Payload.UserID.String() != authorized.User.ID.String() {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched user id")
	}

	if tokenPayload.Payload.AllowIPCode != req.GetCode() {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched code")
	}

	authorized.User.AllowedIps = append(authorized.User.AllowedIps, tokenPayload.Payload.AllowIP)

	err = s.store.UpdateUserAllowIP(ctx, sqlc.UpdateUserAllowIPParams{
		AllowipToken: "null",
		ID:           authorized.User.ID,
		AllowedIps:   authorized.User.AllowedIps,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update allow ip list: %s", err.Error())
	}

	return &rpc.AllowIPResponse{Message: "your ip address has been successfully added to whitelist"}, nil
}

func validateAllowIPRequest(req *rpc.AllowIPRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if len(req.GetCode()) != 6 {
		violations = append(violations, fieldViolation("code", fmt.Errorf("invalid code")))
	}
	return
}
