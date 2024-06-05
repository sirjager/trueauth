package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/sirjager/gopkg/mail"
	"github.com/sirjager/trueauth/config"
	"github.com/sirjager/trueauth/db/db"
)

func RunTaskProcessor(
	ctx context.Context,
	wg *errgroup.Group,
	logr zerolog.Logger,
	store db.Store,
	mailer mail.Sender,
	config config.Config,
	redisOpt asynq.RedisClientOpt,
) {
	processor, err := NewRedisTaskProcessor(logr, store, mailer, config, redisOpt)
	if err != nil {
		logr.Fatal().Err(err).Msg("failed to create task processor")
	}

	logr.Info().Msgf("started task processor")

	if err := processor.Start(); err != nil {
		logr.Fatal().Err(err).Msg("failed to start task processor")
	}

	wg.Go(func() error {
		<-ctx.Done()
		logr.Info().Msg("gracefully shutting down task processor...")
		processor.Shutdown()
		logr.Info().Msg("task process has been stopped")
		return nil
	})
}
