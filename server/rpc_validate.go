package server

import (
	"bytes"
	"context"
	"errors"

	"google.golang.org/grpc/status"

	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/pkg/cache"
	"github.com/sirjager/trueauth/pkg/tokens"
	"github.com/sirjager/trueauth/rpc"
)

func (s *Server) Validate(
	ctx context.Context,
	req *rpc.ValidateRequest,
) (*rpc.ValidateResponse, error) {
	incoming, err := s.tokens.VerifyToken(req.GetToken())
	if err != nil {
		return &rpc.ValidateResponse{Message: err.Error()}, nil
	}

	var stored tokens.Payload
	_key := tokenKey(incoming.Payload.UserID, incoming.Payload.SessionID, TokenTypeAccess)
	if err = s.cache.Get(ctx, _key, &stored); err != nil {
		if errors.Is(cache.ErrNoRecord, err) {
			return &rpc.ValidateResponse{Message: tokens.ErrInvalidToken.Error()}, nil
		}
		return nil, status.Errorf(_internal, err.Error())
	}

	if incoming.ID != stored.ID {
		return &rpc.ValidateResponse{Message: tokens.ErrInvalidToken.Error()}, nil
	}
	if incoming.Payload.SessionID != stored.Payload.SessionID {
		return &rpc.ValidateResponse{Message: tokens.ErrInvalidToken.Error()}, nil
	}
	if !bytes.Equal(incoming.Payload.UserID, stored.Payload.UserID) {
		return &rpc.ValidateResponse{Message: tokens.ErrInvalidToken.Error()}, nil
	}

	user, err := s.store.ReadUser(ctx, incoming.Payload.UserID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return &rpc.ValidateResponse{Message: tokens.ErrInvalidToken.Error()}, nil
		}
		return nil, status.Errorf(_internal, err.Error())
	}

	// verified users are not supposed to have token, though it will never happen
	if !user.Verified {
		return &rpc.ValidateResponse{Message: tokens.ErrInvalidToken.Error()}, nil
	}

	response := &rpc.ValidateResponse{Message: "valid token"}

	if req.GetUser() {
		profile, err := user.Profile()
		if err != nil {
			return nil, status.Errorf(_internal, err.Error())
		}
		response.User = publicProfile(profile)
	}

	return response, nil
}
