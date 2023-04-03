package server

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/sirjager/trueauth/cfg"
	"github.com/sirjager/trueauth/db/sqlc"
	"github.com/sirjager/trueauth/mail"
	"github.com/sirjager/trueauth/worker"
)

func RunTaskProcessor(logger zerolog.Logger, store sqlc.Store, mailer mail.MailSender, config cfg.Config, redisOpt asynq.RedisClientOpt) {
	taskProcessor, err := worker.NewRedisTaskProcessor(logger, store, mailer, config, redisOpt)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create task processor")
	}

	logger.Info().Msgf("started task processor")

	if err := taskProcessor.Start(); err != nil {
		logger.Fatal().Err(err).Msg("failed to start task processor")
	}
}
