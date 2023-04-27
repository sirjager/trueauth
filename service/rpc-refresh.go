package service

import (
	"context"
	"database/sql"
	"errors"

	rpc "github.com/sirjager/rpcs/trueauth/go"
	"github.com/sirjager/trueauth/tokens"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *TrueAuthService) RefreshToken(ctx context.Context, req *rpc.RefreshTokenRequest) (*rpc.RefreshTokenResponse, error) {
	payload, err := s.tokens.VerifyToken(req.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid refresh token: %s", err.Error())
	}
	session, err := s.store.GetSession(ctx, payload.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.Internal, "invalid refresh token")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if session.RefreshToken != req.GetRefreshToken() {
		return nil, status.Errorf(codes.Internal, "invalid refresh token")
	}
	if session.AccountID != payload.Payload.AccountID {
		return nil, status.Errorf(codes.Internal, "invalid refresh token")
	}
	// Generate access tokens
	access, payload, err := s.tokens.CreateToken(tokens.PayloadData{AccountID: session.AccountID}, s.config.AccessTokenTTL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err.Error())
	}
	return &rpc.RefreshTokenResponse{
		AccessToken:          access,
		AccessTokenExpiresAt: timestamppb.New(payload.ExpiresAt),
	}, nil
}
