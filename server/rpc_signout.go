package server

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc/status"

	rpc "github.com/sirjager/trueauth/rpc"
)

func (s *Server) Signout(
	ctx context.Context,
	req *rpc.SignoutRequest,
) (*rpc.SignoutResponse, error) {
	auth, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthorizedError(err)
	}

	removeCookies := func(ctx context.Context) error {
		return s.SetCookies(ctx, []http.Cookie{
			{Name: "sessionId", Value: "", Path: "/", Expires: time.Now(), HttpOnly: true},
			{Name: "accessToken", Value: "", Path: "/", Expires: time.Now(), HttpOnly: true},
			{Name: "refreshToken", Value: "", Path: "/", Expires: time.Now(), HttpOnly: true},
		})
	}

	if req.GetAll() {
		// delete all session for this user
		if err = s.cache.DeleteWithPrefix(ctx, userSessionsKey(auth.User().ID)); err != nil {
			return nil, status.Errorf(_internal, err.Error())
		}
		if err = removeCookies(ctx); err != nil {
			return nil, status.Errorf(_internal, err.Error())
		}
		return &rpc.SignoutResponse{Message: "all sessions deleted"}, nil
	}

	if req.Session != "" {
		targetSession := sessionKey(auth.User().ID, req.Session)
		if err = s.cache.DeleteWithPrefix(ctx, targetSession); err != nil {
			return nil, status.Errorf(_internal, err.Error())
		}
		// if the targeted session is current session then we will remove cookies
		if req.GetSession() == auth.Payload().Payload.SessionID {
			if err = removeCookies(ctx); err != nil {
				return nil, status.Errorf(_internal, err.Error())
			}
		}
		return &rpc.SignoutResponse{Message: "session deleted"}, nil
	}

	currentSession := sessionKey(auth.User().ID, auth.Payload().Payload.SessionID)
	if err = s.cache.DeleteWithPrefix(ctx, currentSession); err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}
	if err = removeCookies(ctx); err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	return &rpc.SignoutResponse{Message: "current session deleted"}, nil
}
