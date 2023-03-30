package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/sirjager/trueauth/cfg"
	"github.com/sirjager/trueauth/cmd/server"
	"github.com/sirjager/trueauth/service"
)

var logger zerolog.Logger
var startTime time.Time
var serviceName string

func init() {
	serviceName = "TrueAuth"
	startTime = time.Now()
	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false})
	logger = logger.With().Timestamp().Logger()
	logger = logger.With().Str("service", strings.ToLower(serviceName)).Logger()
}

func main() {
	config, err := cfg.LoadConfigs(".", "example")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load configurations")
	}
	config.StartTime = startTime
	config.ServiceName = serviceName

	errs := make(chan error)
	go handleSignals(errs)

	srvic, err := service.NewTrueAuthService(logger, config)
	if err != nil {
		logger.Fatal().Err(err).Msgf("error while creating %s service", serviceName)
	}

	if config.RestPort != "" {
		go server.RunGatewayServer(srvic, logger, config, errs)
	}

	if config.GrpcPort != "" {
		go server.RunGRPCServer(srvic, logger, config, errs)
	}

	logger.Error().Err(<-errs).Msg("exit")
}

func handleSignals(errs chan error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errs <- fmt.Errorf("%s", <-c)
}
