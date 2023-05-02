package workers

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/sirjager/trueauth/config"

	"github.com/sirjager/trueauth/internal/db/sqlc"
	"github.com/sirjager/trueauth/internal/worker"

	"github.com/sirjager/trueauth/pkg/mail"
)

func RunTaskProcessor(logger zerolog.Logger, store sqlc.Store, mailer mail.MailSender, config config.Config, redisOpt asynq.RedisClientOpt) {
	taskProcessor, err := worker.NewRedisTaskProcessor(logger, store, mailer, config, redisOpt)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create task processor")
	}

	logger.Info().Msgf("started task processor")

	if err := taskProcessor.Start(); err != nil {
		logger.Fatal().Err(err).Msg("failed to start task processor")
	}
}
