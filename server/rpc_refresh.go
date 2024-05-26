package server

import (
	"context"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/pkg/tokens"
	rpc "github.com/sirjager/trueauth/stubs"
)

const checkRefreshToken = true

func (s *Server) Refresh(
	ctx context.Context,
	req *rpc.RefreshRequest,
) (*rpc.RefreshResponse, error) {
	authorized, err := s.authorize(ctx, checkRefreshToken)
	if err != nil {
		return nil, status.Errorf(_unauthenticated, err.Error())
	}

	meta := s.extractMetadata(ctx)

	tokenParams := tokens.PayloadData{UserID: authorized.User.ID, Type: "access"}
	accessTokenDuration := s.config.Auth.AccessTokenExpDur
	accessToken, accessPayload, err := s.tokens.CreateToken(tokenParams, accessTokenDuration)
	if err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	// NOTE: create new session
	createSessionParams := db.CreateSessionParams{
		Blocked:              false,
		ID:                   accessPayload.ID,
		ClientIp:             meta.ClientIP,
		UserAgent:            meta.UserAgent,
		AccessToken:          accessToken,
		AccessTokenID:        accessPayload.ID,
		AccessTokenExpiresAt: accessPayload.ExpiresAt,
		UserID:               authorized.User.ID,
	}

	// save new session in store
	session, err := s.store.CreateSession(ctx, createSessionParams)
	if err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	response := &rpc.RefreshResponse{
		SessionId:            session.ID,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiresAt),
	}

	if req.GetUser() {
		response.User = authorized.Profile
	}

	return response, nil
}
