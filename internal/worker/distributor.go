package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

type TaskDistributor interface {
	DistributeTaskSendEmailVerified(ctx context.Context, payload PayloadSendEmailVerified, opts ...asynq.Option) error
	DistributeTaskSendEmailVerification(ctx context.Context, payload PayloadSendEmailVerification, opts ...asynq.Option) (err error)
	DistributeTaskClearCompletedVerifications(ctx context.Context, payload PayloadClearCompletedVerfications, opts ...asynq.Option) (err error)
	DistributeTaskSendEmailAllowIP(ctx context.Context, payload PayloadSendEmailAllowIP, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
	logr   zerolog.Logger
}

func NewRedisTaskDistributor(logr zerolog.Logger, redisOpts asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpts)
	return &RedisTaskDistributor{
		client: client,
		logr:   logr,
	}
}
