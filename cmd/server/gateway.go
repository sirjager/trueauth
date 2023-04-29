package server

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/sirjager/trueauth/cfg"
	"github.com/sirjager/trueauth/service"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/sirjager/trueauth/docs/statik"

	rpc "github.com/sirjager/rpcs/trueauth/go"
)

func RunGatewayServer(srvic *service.TrueAuthService, logger zerolog.Logger, config cfg.Config, errs chan error) {
	opts := []runtime.ServeMuxOption{}

	opts = append(opts, runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions:   protojson.MarshalOptions{UseProtoNames: true},
		UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: false},
	}))

	opts = append(opts, runtime.WithIncomingHeaderMatcher(AllowedHeaders([]string{})))

	grpcMux := runtime.NewServeMux(opts...)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := rpc.RegisterTrueAuthHandlerServer(ctx, grpcMux, srvic)
	if err != nil {
		errs <- err
		logger.Fatal().Err(err).Msg("can not register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// File server for swagger documentations
	statikFS, err := fs.New()
	if err != nil {
		errs <- err
		logger.Fatal().Err(err).Msg("can not statik file server")
	}
	swaggerHander := http.StripPrefix("/api/docs/", http.FileServer(statikFS))
	mux.Handle("/api/docs/", swaggerHander)

	mux.Handle("/metrics", promhttp.Handler())

	listener, err := net.Listen("tcp", ":"+config.RestPort)
	if err != nil {
		errs <- err
		logger.Fatal().Err(err).Msg("unable to start rest gateway server")
	}

	logger.Info().Msgf("started REST server at %s", listener.Addr().String())

	handler := HTTPLogger(logger, mux)

	errs <- http.Serve(listener, handler)
}
