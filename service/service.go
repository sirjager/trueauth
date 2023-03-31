package service

import (
	"github.com/rs/zerolog"
	rpc "github.com/sirjager/rpcs/trueauth/go"
	"github.com/sirjager/trueauth/cfg"
	"github.com/sirjager/trueauth/db/sqlc"
)

type TrueAuthService struct {
	rpc.UnimplementedTrueAuthServer
	logger zerolog.Logger
	config cfg.Config
	store  sqlc.Store
}

func NewTrueAuthService(logger zerolog.Logger, config cfg.Config, store sqlc.Store) (*TrueAuthService, error) {
	return &TrueAuthService{
		logger: logger,
		config: config,
		store:  store,
	}, nil
}
