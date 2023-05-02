package gateway

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/sirjager/trueauth/docs/statik"

	"github.com/sirjager/trueauth/internal/service"

	rpc "github.com/sirjager/rpcs/trueauth/go"
)

func RunServer(srvic *service.CoreService, errs chan error) {
	opts := []runtime.ServeMuxOption{}

	opts = append(opts, runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions:   protojson.MarshalOptions{UseProtoNames: true},
		UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: false},
	}))

	allowedHeaders := []string{
		//
	}

	opts = append(opts, runtime.WithIncomingHeaderMatcher(AllowedHeaders(allowedHeaders)))

	grpcMux := runtime.NewServeMux(opts...)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := rpc.RegisterTrueAuthHandlerServer(ctx, grpcMux, srvic)
	if err != nil {
		errs <- err
		srvic.Logr.Fatal().Err(err).Msg("can not register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// File server for swagger documentations
	statikFS, err := fs.New()
	if err != nil {
		errs <- err
		srvic.Logr.Fatal().Err(err).Msg("can not statik file server")
	}
	swaggerHander := http.StripPrefix("/api/docs/", http.FileServer(statikFS))
	mux.Handle("/api/docs/", swaggerHander)

	mux.Handle("/metrics", promhttp.Handler())

	listener, err := net.Listen("tcp", ":"+srvic.Config.GatewayPort)
	if err != nil {
		errs <- err
		srvic.Logr.Fatal().Err(err).Msg("unable to start rest gateway server")
	}

	srvic.Logr.Info().Msgf("started rest server at %s", listener.Addr().String())

	handler := Logger(srvic.Logr, mux)

	errs <- http.Serve(listener, handler)
}
