package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	rpc "github.com/sirjager/rpcs/trueauth/go"

	"github.com/sirjager/trueauth/internal/db/sqlc"

	"github.com/sirjager/trueauth/pkg/crypt"
	"github.com/sirjager/trueauth/pkg/mail"
	"github.com/sirjager/trueauth/pkg/tokens"
	"github.com/sirjager/trueauth/pkg/utils"
	"github.com/sirjager/trueauth/pkg/validator"
)

func (s *CoreService) Login(ctx context.Context, req *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	findBy, violations := validateLoginRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}
	var user sqlc.User
	var err error
	switch strings.ToLower(findBy) {
	case "email":
		user, err = s.store.ReadUserByEmail(ctx, req.GetIdentity())
	default:
		user, err = s.store.ReadUserByUsername(ctx, req.GetIdentity())
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid %s or password", findBy)
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// verify password
	if err := crypt.VerifyPassword(req.GetPassword(), user.Password); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid %s or password", findBy)
	}

	// extract metadata like client-ip and user-agent
	meta := s.extractMetadata(ctx)

	if s.isUnKnownIP(ctx, user) {

		if !user.EmailVerified {
			return nil, status.Errorf(codes.Unauthenticated, "login from unknown ip address with unverified email is not allowed, either login from same ip from where user was created or verify your email address : %s", err.Error())
		}

		// generate a random code
		allowIPCode := utils.RandomNumberAsString(6)
		durationTTL := s.Config.VerifyTokenTTL

		allowIPToken, _, err := s.tokens.CreateToken(
			tokens.PayloadData{
				UserID:      user.ID,
				UserEmail:   user.Email,
				AllowIP:     meta.ClientIp,
				AllowIPCode: allowIPCode,
			}, durationTTL,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create code: %s", err.Error())
		}

		email := mail.Mail{To: []string{user.Email}}
		email.Subject = "Login request from unknown ip address"
		email.Body = fmt.Sprintf(`
		Hello <br/>
		Login request from unknown ip address : <b>%s</b> <br/>
		To allow login from this ip adress use this code : <b>%s</b> <br/>
		This code is only valid for %s <br/> <br/>
		Thank You`, meta.ClientIp, allowIPCode, durationTTL.String())

		err = s.mailer.SendMail(email)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to send email: %s", err.Error())
		}
		err = s.store.UpdateUserAllowIPToken(ctx, sqlc.UpdateUserAllowIPTokenParams{ID: user.ID, AllowipToken: allowIPToken})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update allow ip token: %s", err.Error())
		}

		return nil, status.Errorf(codes.Unauthenticated, "mail has been sent to your email address to allow login from your ip address")
	}

	// Generate tokens, Create sessions and return
	access_token, access_payload, err := s.tokens.CreateToken(tokens.PayloadData{UserID: user.ID}, s.Config.AccessTokenTTL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// generate refresh token
	refresh_token, refresh_payload, err := s.tokens.CreateToken(tokens.PayloadData{UserID: user.ID}, s.Config.RefreshTokenTTL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// create new session
	createSessionParams := sqlc.CreateSessionParams{
		ID:                    refresh_payload.Id,
		ClientIp:              meta.ClientIp,
		UserAgent:             meta.UserAgent,
		RefreshToken:          refresh_token,
		RefreshTokenExpiresAt: refresh_payload.ExpiresAt,
		AccessTokenID:         access_payload.Id,
		AccessToken:           access_token,
		AccessTokenExpiresAt:  access_payload.ExpiresAt,
		UserID:                user.ID,
		Blocked:               false,
	}

	session, err := s.store.CreateSession(ctx, createSessionParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &rpc.LoginResponse{
		User:                  publicProfile(user),
		SessionId:             session.ID.String(),
		AccessToken:           access_token,
		RefreshToken:          refresh_token,
		AccessTokenExpiresAt:  timestamppb.New(access_payload.ExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(refresh_payload.ExpiresAt),
	}, nil
}

func validateLoginRequest(req *rpc.LoginRequest) (identity string, violations []*errdetails.BadRequest_FieldViolation) {
	identity = ""
	if strings.Contains(req.GetIdentity(), "@") {
		if err := validator.ValidateEmail(req.GetIdentity()); err != nil {
			violations = append(violations, fieldViolation("identity", err))
		}
		identity = "email"
	} else {
		if err := validator.ValidateUsername(req.GetIdentity()); err != nil {
			violations = append(violations, fieldViolation("identity", err))
		}
		identity = "username"
	}
	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if identity == "" {
		violations = append(violations, fieldViolation("identity", fmt.Errorf("identity must be a valid username or an email")))
	}
	return identity, violations
}
