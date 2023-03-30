package main

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/sirjager/trueauth/cfg"
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
}
