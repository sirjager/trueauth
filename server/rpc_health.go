package server

import (
	"context"
	"time"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	rpc "github.com/sirjager/trueauth/rpc"
)

func (s *Server) Health(ctx context.Context, req *rpc.HealthRequest) (*rpc.HealthResponse, error) {
	return &rpc.HealthResponse{
		Service:   s.config.ServiceName,
		Server:    s.config.ServerName,
		Status:    healthpb.HealthCheckResponse_SERVING.String(),
		Timestamp: timestamppb.Now(),
		Started:   timestamppb.New(s.config.StartTime),
		Uptime:    durationpb.New(time.Since(s.config.StartTime)),
	}, nil
}
