package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/go-redis/redis/v8"
	rpc "github.com/sirjager/trueauth/stubs/go"

	"github.com/sirjager/trueauth/internal/db/sqlc"
	"github.com/sirjager/trueauth/internal/worker"

	"github.com/sirjager/trueauth/pkg/crypt"
	"github.com/sirjager/trueauth/pkg/tokens"
	"github.com/sirjager/trueauth/pkg/utils"
	"github.com/sirjager/trueauth/pkg/validator"
)

//  1. validate login request
//  2. first we find user
//     2.1 if no user then return: invalid credentials error
//  3. verify password
//     3.1 if password verification failed then return: invalid credentials error
//  4. extracting metadata like: ipaddress, userAgent
//  5. check if current ipaddress is present in user.allowed_ips
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
			return nil, status.Errorf(codes.Unauthenticated, invalidCredentialsError(findBy))
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// verify password
	if err := crypt.VerifyPassword(req.GetPassword(), user.Password); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, invalidCredentialsError(findBy))
	}

	// extract metadata like client-ip and user-agent
	meta := s.extractMetadata(ctx)

	if s.isUnKnownIP(ctx, user) {
		if !user.EmailVerified {
			// Though this will never happen, our registeration process verifies user's email even before creating user.
			return nil, status.Errorf(codes.Unauthenticated, ERR_UNVERIFIED_EMAIL_ADDRESS)
		}

		// before sending email we will check if there is any pending request to allow this ip addr
		// we dont want to send email for same ip address again and agin untill or unless no prev req or code expires
		bytes, err := s.redis.Get(ctx, utils.AllowIPKey(user.Email, meta.ClientIP))
		if err == nil {
			var payload worker.PayloadSendEmailAllowIP
			if err = json.Unmarshal(bytes, &payload); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}

			if req.GetAllowIpCode() != "" {
				if payload.Code == req.GetAllowIpCode() && payload.AllowIP == meta.ClientIP {
					//
				}

			}
			tryAfter := time.Until(payload.LastEmailSentAt.Add(s.Config.AllowIPTokenCooldown))
			return nil, status.Errorf(codes.Internal, "previous request is already pending to allow your current ip address, check instructions on your email or try again after %s", tryAfter)

		}

		// if there is no prev request for current ip address then we send email
		if err = s.taskDistributor.DistributeTaskSendEmailAllowIP(ctx, worker.PayloadSendEmailAllowIP{
			Email:     user.Email,
			AllowIP:   meta.ClientIP,
			UserAgent: meta.UserAgent,
			Timestamp: time.Now(),
		}); err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "mail has been sent to your email address to allow login from your ip address")
	}

	sessionID := utils.UUID_XID()

	accessTokenPayload := tokens.PayloadData{SID: sessionID, UserID: string(user.ID), UserEmail: user.Email, ClientIP: meta.ClientIP}
	access_token, access_payload, err := s.tokens.CreateToken(accessTokenPayload, s.Config.AccessTokenTTL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	refreshTokenPayload := tokens.PayloadData{SID: sessionID, UserID: string(user.ID), UserEmail: user.Email, ClientIP: meta.ClientIP}
	refresh_token, refresh_payload, err := s.tokens.CreateToken(refreshTokenPayload, s.Config.RefreshTokenTTL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	accessRecord := Session{
		ID:        access_payload.ID,
		Token:     access_token,
		ClientIP:  meta.ClientIP,
		UserAgent: meta.UserAgent,
		UserID:    access_payload.Data.UserID,
		ExpiresAt: access_payload.ExpiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Blocked:   false,
	}
	refreshRecord := Session{
		ID:        refresh_payload.ID,
		Token:     refresh_token,
		ClientIP:  meta.ClientIP,
		UserAgent: meta.UserAgent,
		UserID:    refresh_payload.Data.UserID,
		ExpiresAt: refresh_payload.ExpiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Blocked:   false,
	}

	session_access_token_data, err := json.Marshal(accessRecord)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	session_refresh_token_data, err := json.Marshal(refreshRecord)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	access_token_key := accessTokenKey(string(user.ID), access_payload.ID)
	refresh_token_key := refreshTokenKey(string(user.ID), refresh_payload.ID)
	commands := []redis.Cmder{
		redis.NewStringCmd(ctx, "SET", access_token_key, session_access_token_data, "EX", s.Config.AccessTokenTTL.Seconds()),
		redis.NewStatusCmd(ctx, "SET", refresh_token_key, session_refresh_token_data, "EX", s.Config.RefreshTokenTTL.Seconds()),
	}
	_, err = s.redis.Transaction(ctx, commands)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &rpc.LoginResponse{
		User:                  publicProfile(user),
		SessionId:             refresh_payload.ID,
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

	if req.GetAllowIpCode() != "" {
		if len(req.GetAllowIpCode()) != 20 {
			// our code is always unqiue 20char long generated by utils.UUID_XID()
			violations = append(violations, fieldViolation("allow_ip_code", fmt.Errorf("invalid allow ip code")))
		}
	}

	if identity == "" {
		violations = append(violations, fieldViolation("identity", fmt.Errorf("identity must be a valid username or an email")))
	}
	return identity, violations
}

func refreshTokenKey(userid, tokenid string) string {
	return fmt.Sprintf("sessions:%s:refresh:%s", userid, tokenid)
}

func accessTokenKey(userid, tokenid string) string {
	return fmt.Sprintf("sessions:%s:access:%s", userid, tokenid)
}

func invalidCredentialsError(findby string) string {
	return fmt.Sprintf("invalid %s or password", findby)
}

const (
	ERR_UNVERIFIED_EMAIL_ADDRESS = "UNVERIFIED_EMAIL_ADDRESS"
)
