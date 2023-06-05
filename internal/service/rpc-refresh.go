package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirjager/trueauth/pkg/tokens"
	"github.com/sirjager/trueauth/pkg/utils"
	rpc "github.com/sirjager/trueauth/stubs/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	REFRESH_TOKEN_HEADER = "Refresh-Token"
)

func (s *CoreService) Refresh(ctx context.Context, req *rpc.RefreshRequest) (*rpc.RefreshResponse, error) {
	// we are allowing user to send refresh token in two places: headers and body
	// first we check if Refresh-Token header is provided or not, if provided we will use that
	// if header is not provided than we will use it from body,
	// no refresh token not found any where then we return error
	refresh_token := req.GetRefreshToken()

	// we will priorities refresh token in header
	foundTokens, _ := s.extractHeaders(ctx, REFRESH_TOKEN_HEADER)
	if len(foundTokens) > 1 {
		return nil, status.Errorf(codes.InvalidArgument, "%d %s headers found", len(foundTokens), REFRESH_TOKEN_HEADER)
	}

	if len(foundTokens) == 1 {
		refresh_token = foundTokens[0]
	}

	if refresh_token == "" {
		return nil, status.Errorf(codes.InvalidArgument, "refresh token not found, either include in body or in headers")
	}

	payload, err := s.tokens.VerifyToken(refresh_token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	refresh_token_key := refreshTokenKey(payload.Data.UserID, payload.ID)
	sessionBytes, err := s.redis.Get(ctx, refresh_token_key)
	if err != nil {
		if err == redis.Nil {
			return nil, status.Errorf(codes.Internal, "token is invalid")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var stored_session Session
	if err := json.Unmarshal(sessionBytes, &stored_session); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	meta := s.extractMetadata(ctx)

	// match userid
	if payload.Data.UserID != stored_session.UserID {
		return nil, status.Errorf(codes.Internal, "token in invalid")
	}

	// match token string
	if refresh_token != stored_session.Token {
		return nil, status.Errorf(codes.Internal, "token in invalid")
	}

	if stored_session.Blocked {
		return nil, status.Errorf(codes.Internal, "token is blocked")
	}

	// This ensures that refresh token can only be renewed only by the ip address it was assigned to
	if payload.Data.ClientIP != meta.ClientIP || payload.Data.ClientIP != stored_session.ClientIP {
		return nil, status.Errorf(codes.PermissionDenied, "token was not assigned to your current ip address")
	}

	// we will add session ID so that user can logout / delete this session later
	sessionID := utils.UUID_XID()
	accessTokenPayload := tokens.PayloadData{SID: sessionID, UserID: string(payload.Data.UserID), UserEmail: payload.Data.UserEmail, ClientIP: meta.ClientIP}
	access_token, access_payload, err := s.tokens.CreateToken(accessTokenPayload, s.Config.AccessTokenTTL)
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

	access_token_key := accessTokenKey(string(payload.Data.UserID), access_payload.ID)
	err = s.redis.Set(ctx, accessRecord, s.Config.AccessTokenTTL, access_token_key)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &rpc.RefreshResponse{
		SessionId:            sessionID,
		AccessToken:          access_token,
		AccessTokenExpiresAt: timestamppb.New(access_payload.ExpiresAt),
	}, nil
}
