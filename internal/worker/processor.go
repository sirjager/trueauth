package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"

	"github.com/sirjager/trueauth/config"

	"github.com/sirjager/trueauth/internal/db/sqlc"

	"github.com/sirjager/trueauth/pkg/mail"
	"github.com/sirjager/trueauth/pkg/tokens"
)

const (
	QUEUE_CRITICAL = "critical"
	QUEUE_DEFAULT  = "default"
	QUEUE_LOW      = "low"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendEmailVerified(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	logger zerolog.Logger
	store  sqlc.Store
	mailer mail.MailSender
	config config.Config
	tokens tokens.TokenBuilder
}

func NewRedisTaskProcessor(logger zerolog.Logger, store sqlc.Store, mailer mail.MailSender, config config.Config, redisOpts asynq.RedisClientOpt) (TaskProcessor, error) {
	clientConfig := asynq.Config{
		Queues: map[string]int{
			QUEUE_CRITICAL: 10,
			QUEUE_DEFAULT:  5,
			QUEUE_LOW:      1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			logger.Error().Err(err).
				Str("type", task.Type()).
				Interface("payload", task.Payload()).
				Msg("failed to process task")
		}),
		Logger: NewLogger(logger),
	}

	builder, err := tokens.NewPasetoBuilder(config.TokenSecret)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create token builder")
		return nil, err
	}

	server := asynq.NewServer(redisOpts, clientConfig)

	return &RedisTaskProcessor{
		server: server,
		logger: logger,
		store:  store,
		mailer: mailer,
		config: config,
		tokens: builder,
	}, nil
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TASK_SEND_EMAIL_VERIFIED, processor.ProcessTaskSendEmailVerified)

	return processor.server.Start(mux)
}
