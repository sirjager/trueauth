package grpc

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/sirjager/trueauth/rpc"
	"github.com/sirjager/trueauth/server"
)

func RunServer(ctx context.Context, wg *errgroup.Group, address string, srvr *server.Server) {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				Logger(srvr.Logr),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				StreamLogger(srvr.Logr),
			),
		),

		grpc.MaxRecvMsgSize(1024*1024), // bytes * Kilobytes * Megabytes
		grpc.MaxConcurrentStreams(100), // to limit max concurrent streams
	)

	rpc.RegisterTrueAuthServer(grpcServer, srvr)
	reflection.Register(grpcServer)

	// health server to check if server is up
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		srvr.Logr.Fatal().Err(err).Msg("unable to create grpc listener")
	}

	wg.Go(func() error {
		srvr.Logr.Info().Msgf("started grpc server at %s", listener.Addr().String())
		if err = grpcServer.Serve(listener); err != nil {
			srvr.Logr.Error().Err(err).Msg("unable to serve gRPC server")
			return err
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		srvr.Logr.Info().Msg("gracefully shutting down gRPC server")
		grpcServer.GracefulStop()
		srvr.Logr.Info().Msg("gRPC server shutdown complete")
		return nil
	})
}
