package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"

	"github.com/sirjager/trueauth/config"

	"github.com/sirjager/trueauth/internal/db/sqlc"

	"github.com/sirjager/trueauth/pkg/db"
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
	ProcessTaskSendEmailVerification(ctx context.Context, task *asynq.Task) (err error)
	ProcessTaskClearCompletedVerifications(ctx context.Context, task *asynq.Task) (err error)
	ProcessTaskSendEmailAllowIP(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server        *asynq.Server
	logr          zerolog.Logger
	store         sqlc.Store
	mailer        mail.MailSender
	redis         *db.RedisClient
	config        config.Config
	tokens        tokens.TokenBuilder
	emailTemplate string
}

func NewRedisTaskProcessor(
	logger zerolog.Logger,
	store sqlc.Store,
	redis *db.RedisClient,
	mailer mail.MailSender,
	config config.Config,
	redisOpts asynq.RedisClientOpt,
	emailTemplate string,
) (TaskProcessor, error) {
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
		server:        server,
		logr:          logger,
		store:         store,
		redis:         redis,
		mailer:        mailer,
		config:        config,
		tokens:        builder,
		emailTemplate: emailTemplate,
	}, nil
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TASK_SEND_EMAIL_VERIFIED, processor.ProcessTaskSendEmailVerified)
	mux.HandleFunc(TASK_SEND_EMAIL_VERIFICATION, processor.ProcessTaskSendEmailVerification)
	mux.HandleFunc(TASK_CLEAR_COMPLETED_VERIFICATIONS, processor.ProcessTaskClearCompletedVerifications)
	mux.HandleFunc(TASK_SEND_EMAIL_ALLOW_IP, processor.ProcessTaskSendEmailAllowIP)

	return processor.server.Start(mux)
}

const EMAIL_TEMPLATE_HEADING = "{{HEADING}}"
const EMAIL_TEMPLATE_SUBHEADING = "{{SUBHEADING}}"
const EMAIL_TEMPLATE_BODY = "{{BODY}}"
const EMAIL_TEMPLATE_ACTION = "{{ACTION}}"
const EMAIL_TEMPLATE_TEAMNAME = "{{TEAMNAME}}"
