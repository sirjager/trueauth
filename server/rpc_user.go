package server

import (
	"context"

	rpc "github.com/sirjager/trueauth/stubs"
)

// Welcome is a service method that returns a welcome message.
func (s *Server) User(ctx context.Context, r *rpc.UserRequest) (*rpc.UserResponse, error) {
	authorized, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthorizedError(err)
	}

	return &rpc.UserResponse{User: authorized.Profile}, nil
}
