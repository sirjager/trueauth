package server

import (
	"context"
	"fmt"

	rpc "github.com/sirjager/trueauth/stubs"
)

// welcomeMessage returns a formatted welcome message string.
func welcomeMessage(name string) string {
	return fmt.Sprintf("Welcome to %s API", name)
}

// Welcome is a service method that returns a welcome message.
func (s *Server) Welcome(ctx context.Context, r *rpc.WelcomeRequest) (*rpc.WelcomeResponse, error) {
	return &rpc.WelcomeResponse{
		Message: welcomeMessage(s.config.Server.AppName),
		Docs:    "/v1/docs",
	}, nil
}
