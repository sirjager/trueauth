package service

import (
	"github.com/rs/zerolog"
	rpc "github.com/sirjager/rpcs/trueauth/go"

	"github.com/sirjager/trueauth/config"

	"github.com/sirjager/trueauth/internal/db/sqlc"
	"github.com/sirjager/trueauth/internal/worker"

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
}

func NewTrueAuthService(Logr zerolog.Logger, config config.Config, store sqlc.Store, mailer mail.MailSender, taskDistributor worker.TaskDistributor) (*CoreService, error) {
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
	}, nil
}
