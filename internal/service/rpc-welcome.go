package service

import (
	"context"
	"fmt"

	rpc "github.com/sirjager/rpcs/trueauth/go"
)

func welcomeMessage(name string) string {
	return fmt.Sprintf("Welcome to %s Api", name)
}

func (s *CoreService) Welcome(ctx context.Context, req *rpc.WelcomeRequest) (*rpc.WelcomeResponse, error) {
	return &rpc.WelcomeResponse{Message: welcomeMessage(s.Config.ServiceName)}, nil
}
