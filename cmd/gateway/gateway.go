package gateway

import (
	"context"
	"errors"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/sirjager/trueauth/config"
	"github.com/sirjager/trueauth/rpc"
	"github.com/sirjager/trueauth/server"
	_ "github.com/sirjager/trueauth/statik"
)

func StartServer(
	ctx context.Context,
	wg *errgroup.Group,
	address string,
	srvr *server.Server,
	config config.Config,
) {
	// NOTE: Filter headers: + for allowing and - for disallowing
	incomingHeaders := []string{
		"+cookie", //
	}
	outgoingHeaders := []string{
		"+x-server-name",
		"+x-latency",
		"+x-http-code",
	}

	opts := []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   false,
				UseEnumNumbers:  false,
				EmitUnpopulated: false,
			},
		}),
		runtime.WithIncomingHeaderMatcher(customHeaderMatcher(incomingHeaders)),
		runtime.WithForwardResponseOption(addCustomHeaders(config)),
		runtime.WithOutgoingHeaderMatcher(customHeaderMatcher(outgoingHeaders)),
	}

	grpcMux := runtime.NewServeMux(opts...)

	err := rpc.RegisterTrueAuthHandlerServer(ctx, grpcMux, srvr)
	if err != nil {
		srvr.Logr.Fatal().Err(err).Msg("can not register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// File server for swagger documentations
	statikFS, err := fs.New()
	if err != nil {
		srvr.Logr.Fatal().Err(err).Msg("can not statik file server")
	}
	swaggerHander := http.StripPrefix("/docs/", http.FileServer(statikFS))
	mux.Handle("/docs/", swaggerHander)

	handler := logger(srvr.Logr, mux)

	httpServer := &http.Server{Handler: handler, Addr: address}
	wg.Go(func() error {
		srvr.Logr.Info().Msgf("started http server at %s", address)
		if err := httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				srvr.Logr.Error().Err(err).Msg("failed to start http server")
				return err
			}
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		srvr.Logr.Info().Msg("gracefully shutting down http server")
		// NOTE: here we can limit maximum time for graceful shutdown
		// but for now we do not need it, we can use context.Background()
		if err := httpServer.Shutdown(context.Background()); err != nil {
			srvr.Logr.Error().Err(err).Msg("failed to shutdown http server")
			return err
		}
		srvr.Logr.Info().Msg("http server shutdown complete")
		return nil
	})
}
