package server

import (
	"context"

	"github.com/sirjager/trueauth/rpc"
)

// Welcome is a service method that returns a welcome message.
func (s *Server) User(ctx context.Context, r *rpc.UserRequest) (*rpc.UserResponse, error) {
	auth, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthorizedError(err)
	}
	return &rpc.UserResponse{User: auth.Profile()}, nil
}
