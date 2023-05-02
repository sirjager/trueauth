package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

type TaskDistributor interface {
	DistributeTaskSendEmailVerified(ctx context.Context, payload PayloadSendEmailVerified, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
	logger zerolog.Logger
}

func NewRedisTaskDistributor(logger zerolog.Logger, redisOpts asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpts)
	return &RedisTaskDistributor{
		client: client,
		logger: logger,
	}
}
