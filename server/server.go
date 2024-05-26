package server

import (
	"github.com/rs/zerolog"

	"github.com/sirjager/trueauth/config"
	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/pkg/cache"
	"github.com/sirjager/trueauth/pkg/hash"
	"github.com/sirjager/trueauth/pkg/mail"
	"github.com/sirjager/trueauth/pkg/tokens"
	"github.com/sirjager/trueauth/stubs"
	"github.com/sirjager/trueauth/worker"
)

// Server represents the core service of your application.
type Server struct {
	stubs.UnimplementedTrueAuthServer
	Logr   zerolog.Logger
	store  db.Store
	tokens tokens.TokenBuilder
	tasks  worker.TaskDistributor
	mailer mail.Sender
	hasher hash.HashFunction
	cache  cache.Cache
	config config.Config
}

type Adapters struct {
	Logr   zerolog.Logger
	Store  db.Store
	Tasks  worker.TaskDistributor
	Mail   mail.Sender
	Hash   hash.HashFunction
	Cache  cache.Cache
	Tokens tokens.TokenBuilder
	Config config.Config
}

// NewServer creates a new instance of the Service.
func NewServer(adapters *Adapters) *Server {
	return &Server{
		Logr:   adapters.Logr,
		store:  adapters.Store,
		tasks:  adapters.Tasks,
		mailer: adapters.Mail,
		hasher: adapters.Hash,
		cache:  adapters.Cache,
		tokens: adapters.Tokens,
		config: adapters.Config,
	}
}
