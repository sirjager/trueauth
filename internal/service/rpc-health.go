package service

import (
	"context"
	"time"

	rpc "github.com/sirjager/rpcs/trueauth/go"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	StatusUP = "UP"
)

func (s *CoreService) Health(ctx context.Context, req *rpc.HealthRequest) (*rpc.HealthResponse, error) {
	return &rpc.HealthResponse{
		Status:    StatusUP,
		Timestamp: timestamppb.Now(),
		Started:   timestamppb.New(s.Config.StartTime),
		Uptime:    durationpb.New(time.Since(s.Config.StartTime)),
	}, nil
}
