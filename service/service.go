package service

import (
	"github.com/rs/zerolog"
	rpc "github.com/sirjager/rpcs/trueauth/go"
	"github.com/sirjager/trueauth/cfg"
	"github.com/sirjager/trueauth/db/sqlc"
	"github.com/sirjager/trueauth/mail"
	"github.com/sirjager/trueauth/tokens"
	"github.com/sirjager/trueauth/worker"
)

type TrueAuthService struct {
	rpc.UnimplementedTrueAuthServer
	logger          zerolog.Logger
	config          cfg.Config
	store           sqlc.Store
	taskDistributor worker.TaskDistributor
	tokens          tokens.TokenBuilder
	mailer          mail.MailSender
}

func NewTrueAuthService(logger zerolog.Logger, config cfg.Config, store sqlc.Store, mailer mail.MailSender, taskDistributor worker.TaskDistributor) (*TrueAuthService, error) {
	builder, err := tokens.NewPasetoBuilder(config.TokenSecret)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create token builder")
		return nil, err
	}

	return &TrueAuthService{
		logger:          logger,
		config:          config,
		store:           store,
		taskDistributor: taskDistributor,
		tokens:          builder,
		mailer:          mailer,
	}, nil
}
