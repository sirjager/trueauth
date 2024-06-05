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
	"github.com/rs/zerolog/log"
	"github.com/sirjager/gopkg/cache"
	dbPkg "github.com/sirjager/gopkg/db"
	"github.com/sirjager/gopkg/hash"
	"github.com/sirjager/gopkg/mail"
	"github.com/sirjager/gopkg/tokens"
	"github.com/swaggo/swag/example/basic/docs"
	"golang.org/x/sync/errgroup"

	"github.com/sirjager/trueauth/cmd/gateway"
	"github.com/sirjager/trueauth/cmd/grpc"
	"github.com/sirjager/trueauth/config"
	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/logger"
	"github.com/sirjager/trueauth/migrations"
	"github.com/sirjager/trueauth/server"
	"github.com/sirjager/trueauth/worker"
)

var startTime time.Time

// NOTE: Listenting to thse signals for gracefull shutdown
var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func init() {
	startTime = time.Now()
}

func main() {
	// NOTE: change name of .env file here. For defaults, use "defaults"
	config, err := config.LoadConfigs(".", "prod")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configurations")
	}
	config.Server.StartTime = startTime
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)

	logger, err := logger.NewLogger(config.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize logger")
	}
	defer logger.Close()

	ctx, stop := signal.NotifyContext(context.Background(), interuptSignals...)
	defer stop()

	wg, ctx := errgroup.WithContext(ctx)

	// initialize database
	database, conn, err := dbPkg.NewDatabae(ctx, config.Database, logger.Logr)
	if err != nil {
		logger.Logr.Fatal().Err(err).Msg("failed to create new database instance")
	}
	defer database.Close()

	// migrate database to latest version
	if err = database.MigrateUsingBindata(migrations.AssetNames(), migrations.Asset); err != nil {
		logger.Logr.Fatal().Err(err).Msg("failed to migrate database")
	}

	// initialize mailer for sending emails
	mailer, err := mail.NewGmailSender(config.Mail)
	if err != nil {
		logger.Logr.Fatal().Err(err).Msg("failed to initialize gmail smtp")
	}

	// initialize store for database operations
	store := db.NewStore(conn)

	// initialize redis for task distributor
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
	tasks := worker.NewRedisTaskDistributor(logger.Logr, redisOpt)
	defer tasks.Shutdown()

	// redis client for cache system
	redisClient := redis.NewClient(&redis.Options{Addr: redisOpt.Addr})
	if pingErr := redisClient.Ping(ctx).Err(); pingErr != nil {
		logger.Logr.Fatal().Err(pingErr).Msg("failed to ping redis client")
	}
	defer redisClient.Close()
	cache := cache.NewCacheRedis(redisClient, logger.Logr)

	// initialize token builder for token generation
	tokens, err := tokens.NewPasetoBuilder(config.Auth.Secret)
	if err != nil {
		logger.Logr.Fatal().Err(err).Msg("failed to create token builder")
	}

	hasher := hash.NewBryptHash()

	adapters := &server.Adapters{
		Cache:  cache,
		Logr:   logger.Logr,
		Store:  store,
		Tasks:  tasks,
		Mail:   mailer,
		Hash:   hasher,
		Tokens: tokens,
		Config: config,
	}

	// initialize server
	srvr, err := server.New(adapters)
	if err != nil {
		logger.Logr.Fatal().Err(err).Msg("failed to initialize server")
	}

	// start task processor to process tasks in background
	worker.RunTaskProcessor(ctx, wg, logger.Logr, store, mailer, config, redisOpt)

	// start rest server if port is not empty
	if config.Server.RestPort != "" {
		address := config.Server.Host + ":" + config.Server.RestPort
		gateway.StartServer(ctx, wg, address, srvr, config)
	}

	// start grpc server if port is not empty
	if config.Server.GrpcPort != "" {
		address := config.Server.Host + ":" + config.Server.GrpcPort
		grpc.RunServer(ctx, wg, address, srvr)
	}

	err = wg.Wait()
	if err != nil {
		logger.Logr.Fatal().Err(err).Msg("error from wait group")
	}
}
