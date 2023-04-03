package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	rpc "github.com/sirjager/rpcs/trueauth/go"
	"github.com/sirjager/trueauth/db/sqlc"
	"github.com/sirjager/trueauth/mail"
	"github.com/sirjager/trueauth/tokens"
	"github.com/sirjager/trueauth/utils"
	"github.com/sirjager/trueauth/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *TrueAuthService) Login(ctx context.Context, req *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	findBy, violations := validateLoginRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}

	var account sqlc.Account
	var err error

	switch strings.ToLower(findBy) {
	case "email":
		account, err = s.store.GetAccountByEmail(ctx, req.GetIdentity())
	default:
		account, err = s.store.GetAccountByUsername(ctx, req.GetIdentity())
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid %s or password", findBy)
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// verify password
	if err := utils.VerifyPassword(req.GetPassword(), account.Password); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid %s or password", findBy)
	}

	// extract metadata like client-ip and user-agent
	meta := s.extractMetadata(ctx)

	ip, err := s.store.GetIPByAccountID(ctx, account.ID)
	if err != nil {
		// no need to handle no rows error:  first record is created when creating user
		return nil, status.Errorf(codes.Internal, "something went wrong, please try again")
	}

	if s.isBlockedIP(ip, ctx) {
		return nil, status.Errorf(codes.PermissionDenied, "your ip address is in your blacklist.", findBy)
	}

	if !s.isKnownIP(ip, ctx) {

		emailRecord, err := s.store.GetEmailByEmail(ctx, account.Email)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to fetch email record: %s ", err.Error())
		}

		if !emailRecord.Verified {
			return nil, status.Errorf(codes.Unauthenticated, "login from unknown ip address with unverified email is not allowed, either login from same ip from where account was created or verify your email address", err.Error())
		}

		// generate a random code
		sixDigitCode := utils.RandomNumberAsString(6)
		durationTTL := s.config.VerifyTokenTTL

		token, _, err := s.tokens.CreateToken(tokens.PayloadData{AccountEmail: account.Email, AllowIPCode: sixDigitCode}, durationTTL)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create code: %s", err.Error())
		}

		email := mail.Mail{To: []string{account.Email}}
		email.Subject = "Thank you for joining us. Please confirm your email"
		email.Body = fmt.Sprintf(`
		Hello <br/>
		Login request from unknown ip address : <b>%s</b> <br/>
		To allow login from this ip adress use this code : <b>%s</b> <br/>
		This code is only valid for %s <br/> <br/>
		Thank You`, meta.ClientIp, sixDigitCode, durationTTL.String())

		_, err = s.store.UpdateIPTokenTx(ctx, sqlc.UpdateIPTokenTxParams{
			AccountID: account.ID,
			Token:     token,
			BeforeUpdate: func() error {
				return s.mailer.SendMail(email)
			},
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update ip record: %s", err.Error())
		}

		return nil, status.Errorf(codes.Unauthenticated, "mail has been sent to your email address to allow login from your ip address")
	}

	// Generate tokens, Create sessions and return
	access_token, access_payload, err := s.tokens.CreateToken(tokens.PayloadData{AccountID: account.ID}, s.config.AccessTokenTTL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// generate refresh token
	refresh_token, refresh_payload, err := s.tokens.CreateToken(tokens.PayloadData{AccountID: account.ID}, s.config.RefreshTokenTTL)
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
		AccountID:             account.ID,
		Blocked:               false,
	}

	session, err := s.store.CreateSession(ctx, createSessionParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &rpc.LoginResponse{
		Account:               publicProfile(account),
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
