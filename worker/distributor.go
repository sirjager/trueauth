package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

type TaskDistributor interface {
	Shutdown()

	DistributeTaskSendEmailVerification(
		ctx context.Context,
		payload PayloadEmailVerification,
		opts ...asynq.Option,
	) error

	DistributeTaskSendEmailDeletion(
		ctx context.Context,
		payload PayloadEmailDeletion,
		opts ...asynq.Option,
	) error

	DistributeTaskSendPasswordResetCode(
		ctx context.Context,
		payload PayloadPasswordReset,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
	logger zerolog.Logger
}

func NewRedisTaskDistributor(
	logger zerolog.Logger,
	redisOpts asynq.RedisClientOpt,
) TaskDistributor {
	client := asynq.NewClient(redisOpts)
	return &RedisTaskDistributor{client, logger}
}

// close redis
func (d *RedisTaskDistributor) Shutdown() {
	d.client.Close()
}
