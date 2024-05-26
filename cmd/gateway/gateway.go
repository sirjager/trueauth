package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"


	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/sirjager/trueauth/docs/statik"
	"github.com/sirjager/trueauth/server"
	"github.com/sirjager/trueauth/stubs"
)

func StartServer(ctx context.Context, wg *errgroup.Group, address string, srvr *server.Server) {
	// if need any custom headers, add it here to allow
	allowedIncoming := []string{} // only extends default headers
	allowedOutgoing := []string{}

	opts := []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: false},
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   false,
				UseEnumNumbers:  false,
				EmitUnpopulated: false,
			},
		}),
		runtime.WithIncomingHeaderMatcher(customHeaderMatcher(srvr.Logr, allowedIncoming)),
		runtime.WithForwardResponseOption(mutateResponse(srvr.Logr)),
		runtime.WithOutgoingHeaderMatcher(customHeaderMatcher(srvr.Logr, allowedOutgoing)),
	}

	grpcMux := runtime.NewServeMux(opts...)

	err := stubs.RegisterTrueAuthHandlerServer(ctx, grpcMux, srvr)
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
	swaggerHander := http.StripPrefix("/v1/docs/", http.FileServer(statikFS))
	mux.Handle("/v1/docs/", swaggerHander)

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
		fmt.Println()
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
