package server

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/sirjager/trueauth/db/sqlc"
	"github.com/sirjager/trueauth/worker"
)

func RunTaskProcessor(logger zerolog.Logger, store sqlc.Store, redisOpt asynq.RedisClientOpt) {
	taskProcessor := worker.NewRedisTaskProcessor(logger, store, redisOpt)
	logger.Info().Msgf("started task processor")
	if err := taskProcessor.Start(); err != nil {
		logger.Fatal().Err(err).Msg("failed to start task processor")
	}
}
