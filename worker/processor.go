package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/sirjager/trueauth/db/sqlc"
)

const (
	QUEUE_CRITICAL = "critical"
	QUEUE_DEFAULT  = "default"
	QUEUE_LOW      = "low"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	logger zerolog.Logger
	store  sqlc.Store
}

func NewRedisTaskProcessor(logger zerolog.Logger, store sqlc.Store, redisOpts asynq.RedisClientOpt) TaskProcessor {
	clientConfig := asynq.Config{
		Queues: map[string]int{
			QUEUE_CRITICAL: 10,
			QUEUE_DEFAULT:  5,
			QUEUE_LOW:      1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			logger.Error().Err(err).
				Str("type", task.Type()).
				Bytes("payload", task.Payload()).
				Msg("failed to process task")
		}),
		Logger: NewLogger(logger),
	}
	server := asynq.NewServer(redisOpts, clientConfig)
	return &RedisTaskProcessor{server, logger, store}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TASK_SEND_VERIFY_EMAIL, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
