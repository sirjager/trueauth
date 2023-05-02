package server

import (
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/sirjager/trueauth/cfg"
	"github.com/sirjager/trueauth/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	rpc "github.com/sirjager/rpcs/trueauth/go"
)

func RunGRPCServer(srvic *service.TrueAuthService, logger zerolog.Logger, config cfg.Config, errs chan error) {
	listener, err := net.Listen("tcp", ":"+config.GrpcPort)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to listen grpc tcp server")
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				GRPCLogger(logger),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				GRPCStreamLogger(logger),
			),
		),

		grpc.MaxRecvMsgSize(1024*1024), // bytes * Kilobytes * Megabytes
	)

	rpc.RegisterTrueAuthServer(grpcServer, srvic)

	reflection.Register(grpcServer)

	http.Handle("/metrics", promhttp.Handler())

	logger.Info().Msgf("started gRPC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to serve gRPC server")
	}
}
