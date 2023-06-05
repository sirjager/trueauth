package service

import (
	"github.com/rs/zerolog"
	rpc "github.com/sirjager/trueauth/stubs/go"

	"github.com/sirjager/trueauth/config"

	"github.com/sirjager/trueauth/internal/db/sqlc"
	"github.com/sirjager/trueauth/internal/worker"

	"github.com/sirjager/trueauth/pkg/db"
	"github.com/sirjager/trueauth/pkg/mail"
	"github.com/sirjager/trueauth/pkg/tokens"
)

type CoreService struct {
	rpc.UnimplementedTrueAuthServer
	Logr            zerolog.Logger
	Config          config.Config
	store           sqlc.Store
	tokens          tokens.TokenBuilder
	mailer          mail.MailSender
	taskDistributor worker.TaskDistributor
	redis           *db.RedisClient
}

func NewTrueAuthService(
	Logr zerolog.Logger,
	config config.Config,
	store sqlc.Store,
	mailer mail.MailSender,
	taskDistributor worker.TaskDistributor,
	redisClient *db.RedisClient,
) (*CoreService, error) {
	builder, err := tokens.NewPasetoBuilder(config.TokenSecret)
	if err != nil {
		Logr.Error().Err(err).Msg("failed to create token builder")
		return nil, err
	}

	return &CoreService{
		Logr:            Logr,
		Config:          config,
		store:           store,
		tokens:          builder,
		mailer:          mailer,
		taskDistributor: taskDistributor,
		redis:           redisClient,
	}, nil
}
