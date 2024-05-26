package server

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc/status"

	"github.com/sirjager/trueauth/pkg/utils"
	rpc "github.com/sirjager/trueauth/stubs"
)

func (s *Server) Signout(
	ctx context.Context,
	req *rpc.SignoutRequest,
) (*rpc.SignoutResponse, error) {
	authorized, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthorizedError(err)
	}

	// no matter what the session is we need to clear cookies from client
	if err = s.sendCookies(ctx, []http.Cookie{
		{Name: "sessionId", Value: "", Path: "/", Expires: time.Now(), HttpOnly: true},
		{Name: "accessToken", Value: "", Path: "/", Expires: time.Now(), HttpOnly: true},
		{Name: "refreshToken", Value: "", Path: "/", Expires: time.Now(), HttpOnly: true},
	}); err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	// if all is true, delete all sessions
	if req.GetAll() {
		if err = s.store.DeleteSessionByUserID(ctx, authorized.User.ID); err != nil {
			return nil, status.Errorf(_internal, err.Error())
		}
		return &rpc.SignoutResponse{Message: "all sessions deleted"}, nil
	}

	// if session is not empty, delete targeted session
	if req.GetSession() != "" {
		id, iErr := utils.XIDFromString(req.GetSession())
		if iErr != nil {
			return nil, status.Errorf(_invalidArgument, "invalid session id")
		}
		if err = s.store.DeleteSession(ctx, id.Bytes()); err != nil {
			return nil, status.Errorf(_internal, err.Error())
		}

		return &rpc.SignoutResponse{Message: "session deleted"}, nil
	}

	// if session is empty, delete current session
	// it extracts access tokens, and access token payload if any
	meta := s.extractMetadata(ctx)

	if meta.payload == nil {
		return nil, status.Errorf(_unauthenticated, errUnauthorized)
	}

	// delete session using access token id
	if err := s.store.DeleteSessionByAccessTokenID(ctx, meta.payload.ID); err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	return &rpc.SignoutResponse{Message: "current session deleted"}, nil
}
