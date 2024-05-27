package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/swaggo/swag/example/basic/docs"
	"golang.org/x/sync/errgroup"

	"github.com/sirjager/trueauth/cmd/gateway"
	"github.com/sirjager/trueauth/cmd/grpc"
	"github.com/sirjager/trueauth/config"
	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/pkg/cache"
	dbPkg "github.com/sirjager/trueauth/pkg/db"
	"github.com/sirjager/trueauth/pkg/hash"
	"github.com/sirjager/trueauth/pkg/mail"
	"github.com/sirjager/trueauth/pkg/tokens"
	"github.com/sirjager/trueauth/server"
	"github.com/sirjager/trueauth/worker"
)

var (
	logr      zerolog.Logger
	startTime time.Time
)

// NOTE: Listenting to thse signals for gracefull shutdown
var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func init() {
	startTime = time.Now()
	logr = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false})
	logr = logr.With().Timestamp().Logger()
}

func main() {
	// INFO: change name of .env file here. For defaults, use "defaults"
	config, err := config.LoadConfigs(".", "prod")
	if err != nil {
		logr.Fatal().Err(err).Msg("failed to load configurations")
	}

	ctx, stop := signal.NotifyContext(context.Background(), interuptSignals...)
	defer stop()

	wg, ctx := errgroup.WithContext(ctx)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)

	// NOTE: store server start time in config
	config.Server.StartTime = startTime
	// NOTE: store email template in config

	// NOTE: update server name in logger
	logr = logr.With().Str("server", config.Server.ServerName).Logger()

	// NOTE: initialize database
	database, conn, err := dbPkg.NewDatabae(ctx, config.Database, logr)
	if err != nil {
		logr.Fatal().Err(err).Msg("failed to create new database instance")
	}
	defer database.Close()

	// NOTE: migrate database to latest version
	if err = database.MigrateUsingBindata(); err != nil {
		logr.Fatal().Err(err).Msg("failed to migrate database")
	}

	// NOTE: initialize mailer for sending emails
	mailer, err := mail.NewGmailSender(config.Mail)
	if err != nil {
		logr.Fatal().Err(err).Msg("failed to initialize gmail smtp")
	}

	// NOTE: initialize store for database operations
	store := db.NewStore(conn)

	// NOTE: initialize redis for task distributor
	redisOpt := asynq.RedisClientOpt{
		Addr:     config.Database.RedisAddr,
		Password: config.Database.RedisPass,
		Username: config.Database.RedisUser,
	}
	if config.Database.RedisURL != "" {
		if opts, pErr := redis.ParseURL(config.Database.RedisURL); pErr != nil {
			redisOpt.Addr = opts.Addr
			redisOpt.Username = opts.Username
			redisOpt.Password = opts.Password
			redisOpt.TLSConfig = opts.TLSConfig
			redisOpt.Network = opts.Network
		}
	}
	tasks := worker.NewRedisTaskDistributor(logr, redisOpt)
	defer tasks.Shutdown()

	// NOTE: redis client for cache system
	redisClient := redis.NewClient(&redis.Options{Addr: redisOpt.Addr})
	if pingErr := redisClient.Ping(ctx).Err(); pingErr != nil {
		logr.Fatal().Err(pingErr).Msg("failed to ping redis client")
	}
	defer redisClient.Close()
	cache := cache.NewCacheRedis(redisClient, logr)

	// NOTE: initialize token builder for token generation
	tokens, err := tokens.NewPasetoBuilder(config.Auth.Secret)
	if err != nil {
		logr.Fatal().Err(err).Msg("failed to create token builder")
	}

	hasher := hash.NewBryptHash()

	adapters := &server.Adapters{
		Cache:  cache,
		Logr:   logr,
		Store:  store,
		Tasks:  tasks,
		Mail:   mailer,
		Hash:   hasher,
		Tokens: tokens,
		Config: config,
	}

	// NOTE: initialize server with logr, config, store, tasks, mailer, and tokens
	srvr := server.NewServer(adapters)

	// NOTE: start rest server if port is not empty
	if config.Server.RestPort != "" {
		address := config.Server.Host + ":" + config.Server.RestPort
		gateway.StartServer(ctx, wg, address, srvr)
	}

	// NOTE: start grpc server if port is not empty
	if config.Server.GrpcPort != "" {
		address := config.Server.Host + ":" + config.Server.GrpcPort
		grpc.RunServer(ctx, wg, address, srvr)
	}

	// NOTE: start task processor to process tasks in background
	worker.RunTaskProcessor(ctx, wg, logr, store, mailer, config, redisOpt)

	err = wg.Wait()
	if err != nil {
		logr.Fatal().Err(err).Msg("error from wait group")
	}
}
