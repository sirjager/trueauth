package grpc

import (
	"net"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sirjager/trueauth/internal/service"

	rpc "github.com/sirjager/rpcs/trueauth/go"
)

func RunServer(srvic *service.CoreService, errs chan error) {
	listener, err := net.Listen("tcp", ":"+srvic.Config.GrpcPort)
	if err != nil {
		srvic.Logr.Fatal().Err(err).Msg("unable to listen grpc tcp server")
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				Logger(srvic.Logr),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				StreamLogger(srvic.Logr),
			),
		),

		grpc.MaxRecvMsgSize(1024*1024), // bytes * Kilobytes * Megabytes
	)

	rpc.RegisterTrueAuthServer(grpcServer, srvic)

	reflection.Register(grpcServer)

	http.Handle("/metrics", promhttp.Handler())

	srvic.Logr.Info().Msgf("started grpc server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		srvic.Logr.Fatal().Err(err).Msg("unable to serve gRPC server")
	}
}
