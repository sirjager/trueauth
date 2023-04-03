package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/sirjager/trueauth/cfg"
	"github.com/sirjager/trueauth/cmd/server"
	"github.com/sirjager/trueauth/db"
	"github.com/sirjager/trueauth/db/sqlc"
	"github.com/sirjager/trueauth/mail"
	"github.com/sirjager/trueauth/service"
	"github.com/sirjager/trueauth/worker"
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
	config, err := cfg.LoadConfigs(".", "remote")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load configurations")
	}
	config.StartTime = startTime
	config.ServiceName = serviceName

	conn, err := sql.Open(config.DBConfig.DBDriver, config.DBConfig.DBUrl)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to make database connection")
	}
	defer conn.Close()

	if err = db.PingRedis(config.DBConfig.RedisAddr, logger); err != nil {
		logger.Fatal().Err(err).Msg("failed to ping redis")
	}

	if err = db.Migrate(logger, conn, config.DBConfig); err != nil {
		logger.Fatal().Err(err).Msg("failed to migrate database")
	}

	mailer, err := mail.NewGmailSender(config.GmailSMTP)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize gmail smtp")
	}

	store := sqlc.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{Addr: config.DBConfig.RedisAddr}
	taskDistributor := worker.NewRedisTaskDistributor(logger, redisOpt)
	go server.RunTaskProcessor(logger, store, mailer, config, redisOpt)

	errs := make(chan error)
	go handleSignals(errs)

	srvic, err := service.NewTrueAuthService(logger, config, store, mailer, taskDistributor)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create service")
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
