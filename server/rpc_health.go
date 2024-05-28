package server

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	rpc "github.com/sirjager/trueauth/rpc"
)

const (
	StatusUP = "UP"
)

func (s *Server) Health(ctx context.Context, req *rpc.HealthRequest) (*rpc.HealthResponse, error) {
	s.authorize(ctx)
	return &rpc.HealthResponse{
		Status:    StatusUP,
		Timestamp: timestamppb.Now(),
		Started:   timestamppb.New(s.config.Server.StartTime),
		Uptime:    durationpb.New(time.Since(s.config.Server.StartTime)),
	}, nil
}
