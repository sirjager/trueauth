package service

import (
	"github.com/rs/zerolog"
	rpc "github.com/sirjager/rpcs/trueauth/go"
	"github.com/sirjager/trueauth/cfg"
)

type TrueAuthService struct {
	rpc.UnimplementedTrueAuthServer
	logger zerolog.Logger
	config cfg.Config
}

func NewTrueAuthService(logger zerolog.Logger, config cfg.Config) (*TrueAuthService, error) {
	return &TrueAuthService{
		logger: logger,
		config: config,
	}, nil
}
