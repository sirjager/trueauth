package server

import (
	"context"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sirjager/trueauth/pkg/tokens"
	"github.com/sirjager/trueauth/pkg/utils"
	"github.com/sirjager/trueauth/rpc"
)

const checkRefreshToken = true

func (s *Server) Refresh(
	ctx context.Context,
	req *rpc.RefreshRequest,
) (*rpc.RefreshResponse, error) {
	auth, err := s.authorize(ctx, checkRefreshToken)
	if err != nil {
		return nil, status.Errorf(_unauthenticated, err.Error())
	}

	meta := s.extractMetadata(ctx)

	sessionID := utils.XIDNew().String()
	tokenPayload := tokens.PayloadData{
		UserID:    auth.user.ID,
		ClientIP:  meta.clientIP,
		UserAgent: meta.userAgent,
		SessionID: sessionID,
	}
	accessTokenDuration := s.config.Auth.AccessTokenExpDur
	accessToken, accessPayload, err := s.tokens.CreateToken(tokenPayload, accessTokenDuration)
	if err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	accessKey := tokenKey(auth.user.ID, sessionID, TokenTypeAccess)
	if err = s.cache.Set(ctx, accessKey, accessPayload, accessTokenDuration); err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}
	response := &rpc.RefreshResponse{
		SessionId:            sessionID,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiresAt),
		Message:              "refreshed successfully",
	}

	if req.GetUser() {
		_profile, err := auth.user.Profile()
		if err != nil {
			return nil, status.Errorf(_internal, err.Error())
		}
		response.User = publicProfile(_profile)
	}

	return response, nil
}
