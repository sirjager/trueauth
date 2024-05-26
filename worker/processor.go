package worker

import (
	"context"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"

	"github.com/sirjager/trueauth/config"
	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/pkg/mail"
	"github.com/sirjager/trueauth/pkg/tokens"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskSendEmailVerificationCode(ctx context.Context, task *asynq.Task) error
	ProcessTaskSendUserDeletionCode(ctx context.Context, task *asynq.Task) error
	ProcessTaskSendPasswordResetCode(ctx context.Context, task *asynq.Task) error
}

// RedisTaskProcessor implements the TaskProcessor interface
type RedisTaskProcessor struct {
	logger zerolog.Logger
	store  db.Store
	mailer mail.Sender
	tokens tokens.TokenBuilder
	server *asynq.Server
	config config.Config
}

// NewRedisTaskProcessor creates a new RedisTaskProcessor
func NewRedisTaskProcessor(
	logger zerolog.Logger,
	store db.Store,
	mailer mail.Sender,
	config config.Config,
	redisOpts asynq.RedisClientOpt,
) (TaskProcessor, error) {
	clientConfig := asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
			QueueLow:      1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(
			func(ctx context.Context, task *asynq.Task, err error) {
				logger.Error().Err(err).
					Str("type", task.Type()).
					Interface("payload", task.Payload()).
					Msg("failed to process task")
			},
		),
		Logger: NewLogger(logger),
	}

	builder, err := tokens.NewPasetoBuilder(config.Auth.Secret)
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

// Start starts the RedisTaskProcessor
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendEmailVerificationCode, processor.ProcessTaskSendEmailVerificationCode)
	mux.HandleFunc(TaskSendUserDeletionCode, processor.ProcessTaskSendUserDeletionCode)
	mux.HandleFunc(TaskSendPasswordResetCode, processor.ProcessTaskSendPasswordResetCode)

	return processor.server.Start(mux)
}

func (processor *RedisTaskProcessor) Shutdown() {
	processor.server.Shutdown()
}
