package server

import (
	"context"

	"google.golang.org/grpc/status"

	"github.com/sirjager/trueauth/rpc"
)

// Welcome is a service method that returns a welcome message.
func (s *Server) User(ctx context.Context, r *rpc.UserRequest) (*rpc.UserResponse, error) {
	auth, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthorizedError(err)
	}
	profile, err := auth.user.Profile()
	if err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	return &rpc.UserResponse{User: publicProfile(profile)}, nil
}
