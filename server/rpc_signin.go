package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirjager/gopkg/utils"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sirjager/trueauth/pkg/tokens"
	rpc "github.com/sirjager/trueauth/rpc"
)

func sessionKey(userID []byte, sessionID string) string {
	return fmt.Sprintf("sess:%s:%s", utils.BytesToBase64(userID), sessionID)
}

func userSessionsKey(userID []byte) string {
	return fmt.Sprintf("sess:%s", utils.BytesToBase64(userID))
}

func tokenKey(userID []byte, sessionID, tokenType string) string {
	return fmt.Sprintf(
		"sess:%s:%s:%s", // sess:userID:sessionID:(refresh|access)
		utils.BytesToBase64(userID), sessionID, tokenType,
	)
}

func (s *Server) Signin(
	ctx context.Context,
	req *rpc.SigninRequest,
) (*rpc.SigninResponse, error) {
	auth, err := s.authenticate(ctx)
	if err != nil {
		return nil, status.Errorf(_unauthenticated, err.Error())
	}

	sessionID := utils.XIDNew().String()

	tokenPayload := tokens.PayloadData{
		UserID:    auth.Profile().Id,
		ClientIP:  auth.ClientIP(),
		UserAgent: auth.UserAgent(),
		SessionID: sessionID,
	}

	aTDuration := s.config.Auth.AccessTokenExpDur
	aToken, aPayload, err := s.tokens.CreateToken(tokenPayload, aTDuration)
	if err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	rTDuration := s.config.Auth.RefreshTokenExpDur
	rToken, rPayload, err := s.tokens.CreateToken(tokenPayload, rTDuration)
	if err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	accessKey := tokenKey(auth.Profile().Id, sessionID, TokenTypeAccess)
	if err = s.cache.Set(ctx, accessKey, aPayload, aTDuration); err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	refreshKey := tokenKey(auth.Profile().Id, sessionID, TokenTypeRefresh)
	if err = s.cache.Set(ctx, refreshKey, rPayload, rTDuration); err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	response := &rpc.SigninResponse{Message: "successfully signed in"}

	if req.GetTokens() {
		response.SessionId = sessionID
		response.AccessToken = aToken
		response.RefreshToken = rToken
		response.AccessTokenExpiresAt = timestamppb.New(aPayload.ExpiresAt)
		response.RefreshTokenExpiresAt = timestamppb.New(rPayload.ExpiresAt)
	}

	if req.GetUser() {
		response.User = auth.Profile()
	}

	if req.GetCookies() {
		if err = s.sendCookies(ctx, []http.Cookie{
			{
				Name: "sessionId", Value: sessionID,
				Path: "/", Expires: aPayload.ExpiresAt,
				HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
			},
			{
				Name: "accessToken", Value: aToken,
				Path: "/", Expires: aPayload.ExpiresAt,
				HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
			},
			{
				Name: "refreshToken", Value: rToken,
				Path: "/", Expires: rPayload.ExpiresAt,
				HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
			},
		}); err != nil {
			return nil, status.Errorf(_internal, err.Error())
		}
	}

	return response, nil
}
