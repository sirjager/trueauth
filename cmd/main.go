package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"

	"github.com/sirjager/trueauth/config"

	"github.com/sirjager/trueauth/cmd/gateway"
	"github.com/sirjager/trueauth/cmd/grpc"
	"github.com/sirjager/trueauth/cmd/workers"

	"github.com/sirjager/trueauth/internal/db/sqlc"
	"github.com/sirjager/trueauth/internal/service"
	"github.com/sirjager/trueauth/internal/worker"

	"github.com/sirjager/trueauth/pkg/db"
	"github.com/sirjager/trueauth/pkg/mail"
)

var logr zerolog.Logger
var startTime time.Time
var serviceName string

func init() {
	serviceName = "TrueAuth"
	startTime = time.Now()
	logr = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false})
	logr = logr.With().Timestamp().Logger()
	logr = logr.With().Str("service", strings.ToLower(serviceName)).Logger()
}

func main() {
	config, err := config.LoadConfigs(".", "example")
	if err != nil {
		logr.Fatal().Err(err).Msg("failed to load configurations")
	}
	config.StartTime = startTime
	config.ServiceName = serviceName

	database, conn, err := db.NewDatabae(config.Database, logr)
	if err != nil {
		logr.Fatal().Err(err).Msg("failed to create new database instance")
	}
	defer database.Close()

	if err = database.Ping(); err != nil {
		logr.Fatal().Err(err).Msg("failed to ping database")
	}

	if err = database.Migrate(); err != nil {
		logr.Fatal().Err(err).Msg("failed to migrate database")
	}

	mailer, err := mail.NewGmailSender(config.Mail)
	if err != nil {
		logr.Fatal().Err(err).Msg("failed to initialize gmail smtp")
	}

	store := sqlc.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{Addr: config.Database.RedisAddr}
	taskDistributor := worker.NewRedisTaskDistributor(logr, redisOpt)
	go workers.RunTaskProcessor(logr, store, mailer, config, redisOpt)

	errs := make(chan error)
	go handleSignals(errs)

	srvic, err := service.NewTrueAuthService(logr, config, store, mailer, taskDistributor)
	if err != nil {
		logr.Fatal().Err(err).Msg("failed to create service")
	}

	if config.GatewayPort != "" {
		go gateway.RunServer(srvic, errs)
	}

	if config.GrpcPort != "" {
		go grpc.RunServer(srvic, errs)
	}

	logr.Error().Err(<-errs).Msg("exit")
}

func handleSignals(errs chan error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errs <- fmt.Errorf("%s", <-c)
}
