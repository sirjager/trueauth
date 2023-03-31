package service

import (
	"github.com/rs/zerolog"
	rpc "github.com/sirjager/rpcs/trueauth/go"
	"github.com/sirjager/trueauth/cfg"
	"github.com/sirjager/trueauth/db/sqlc"
	"github.com/sirjager/trueauth/worker"
)

type TrueAuthService struct {
	rpc.UnimplementedTrueAuthServer
	logger          zerolog.Logger
	config          cfg.Config
	store           sqlc.Store
	taskDistributor worker.TaskDistributor
}

func NewTrueAuthService(logger zerolog.Logger, config cfg.Config, store sqlc.Store, taskDistributor worker.TaskDistributor) (*TrueAuthService, error) {
	return &TrueAuthService{
		logger:          logger,
		config:          config,
		store:           store,
		taskDistributor: taskDistributor,
	}, nil
}
