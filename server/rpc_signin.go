package server

import (
	"context"
	"net/http"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/pkg/tokens"
	"github.com/sirjager/trueauth/pkg/utils"
	rpc "github.com/sirjager/trueauth/stubs"
)

const (
	RUNNING = "running"
)

func (s *Server) Signin(ctx context.Context, req *rpc.SigninRequest) (*rpc.SigninResponse, error) {
	s.Logr.Log().Str("method", "Signin").Msg(RUNNING)

	authenticated, err := s.authenticate(ctx)
	if err != nil {
		return nil, status.Errorf(_unauthenticated, err.Error())
	}

	meta := s.extractMetadata(ctx)

	tokenParams := tokens.PayloadData{UserID: authenticated.Profile.ID, Type: "access"}
	accessTokenDuration := s.config.Auth.AccessTokenExpDur
	accessToken, accessPayload, err := s.tokens.CreateToken(tokenParams, accessTokenDuration)
	if err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	tokenParams.Type = "refresh"
	refreshTokenDuration := s.config.Auth.RefreshTokenExpDur
	refreshToken, refreshPayload, err := s.tokens.CreateToken(tokenParams, refreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	// create new session
	createSessionParams := db.CreateSessionParams{
		ID:                    refreshPayload.ID,
		ClientIp:              meta.ClientIP,
		UserAgent:             meta.UserAgent,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt,
		AccessTokenID:         accessPayload.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt,
		UserID:                authenticated.Profile.ID,
		Blocked:               false,
	}

	// save new session in store
	session, err := s.store.CreateSession(ctx, createSessionParams)
	if err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	response := &rpc.SigninResponse{Message: "successfully signed in"}

	if req.GetTokens() {
		response.SessionId = session.ID
		response.AccessToken = accessToken
		response.RefreshToken = refreshToken
		response.AccessTokenExpiresAt = timestamppb.New(accessPayload.ExpiresAt)
		response.RefreshTokenExpiresAt = timestamppb.New(refreshPayload.ExpiresAt)
	}

	if req.GetUser() {
		response.User = publicProfile(authenticated.Profile)
	}

	if req.GetCookies() {
		if err = s.sendCookies(ctx, []http.Cookie{
			{
				Name: "sessionId", Value: utils.BytesToBase64(session.ID),
				Path: "/", Expires: accessPayload.ExpiresAt,
				HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
			},
			{
				Name: "accessToken", Value: accessToken,
				Path: "/", Expires: accessPayload.ExpiresAt,
				HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
			},
			{
				Name: "refreshToken", Value: refreshToken,
				Path: "/", Expires: refreshPayload.ExpiresAt,
				HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
			},
		}); err != nil {
			return nil, status.Errorf(_internal, err.Error())
		}
	}

	return response, nil
}
