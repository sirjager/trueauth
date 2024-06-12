package server

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/sirjager/trueauth/db/db"
	rpc "github.com/sirjager/trueauth/rpc"
)

type Authenticated interface {
	Profile() *rpc.User
	ClientIP() string
	UserAgent() string
}

type _authenticated struct {
	dbuser db.User
	meta   MetaData
}

const (
	errInvalidCredentials = "invalid credentials"
)

// Used for Signing in with username and password
//
// Checks Authorization header for Basic Auth Only, verifies credentials and returns user and error
func (s *Server) authenticate(ctx context.Context) (Authenticated, error) {
	var err error
	// extracting metadata from FromIncomingContext
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = fmt.Errorf("missing metadata from context")
		s.Logr.Error().Err(err).Msg("missing metadata from context")
		return nil, err
	}

	//  extracting authorizationHeader
	values := meta.Get("authorization")
	if len(values) == 0 {
		// return error if no authorization header
		s.Logr.Error().Err(err).Msg("missing authorization header")
		return nil, fmt.Errorf("missing authorization header")
	}

	// checking auth token format
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		// return error if invalid authorization header
		s.Logr.Error().Err(err).Msg("invalid authorization header")
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := fields[0]
	if strings.ToLower(authType) != "basic" {
		// return error if unsupported authorization type
		s.Logr.Error().Err(err).Msg("unsupported authorization type")
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	// storing payload and token in authorized if valid token
	decoded, err := base64.StdEncoding.DecodeString(fields[1])
	if err != nil {
		// return error if failed to decode base64 authorization
		s.Logr.Error().Err(err).Msg("failed to decode base64 authorization")
		return nil, fmt.Errorf("failed to decode base64: %s", err.Error())
	}
	credentials := strings.Split(string(decoded), ":")
	identity := credentials[0]
	password := credentials[1]

	findBy := "username"
	if validateErr := validateEmail(identity); validateErr == nil {
		findBy = "email"
	}

	var dbuser db.User
	switch findBy {
	case "email":
		dbuser, err = s.store.ReadUserByEmail(ctx, identity)
	default:
		dbuser, err = s.store.ReadUserByUsername(ctx, identity)
	}
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			// return error invalid credentials even if user does not exist
			// we don't want to disclose if user exists or not
			s.Logr.Error().Err(err).Msg("user not found") // we can optionally log the real error
			return nil, fmt.Errorf(errInvalidCredentials)
		}
		// if the error is realated to database, return error with real error
		s.Logr.Error().Err(err).Msg("failed to get user")
		return nil, fmt.Errorf(errFailedToRetrieveUser, err.Error())
	}

	if !dbuser.Verified {
		// return error if email not verified
		s.Logr.Error().Err(err).Msg("email not verified")
		return nil, fmt.Errorf(errEmailNotVerified)
	}

	if err = s.hasher.Verify(dbuser.HashSalt, dbuser.HashPass, password); err != nil {
		// return error invalid credentials if password does not match
		s.Logr.Error().Err(err).Msg("invalid password") // log the real error
		return nil, fmt.Errorf(errInvalidCredentials)
	}

	// extracting metadata from FromIncomingContext,
	// so we can use it later in service methods
	_meta := s.extractMetadata(ctx)
	return &_authenticated{
		dbuser: dbuser,
		meta:   _meta,
	}, err
}

func (a *_authenticated) Profile() *rpc.User {
	return publicProfile(a.dbuser)
}

func (a *_authenticated) ClientIP() string {
	return a.meta.clientIP()
}

func (a *_authenticated) UserAgent() string {
	return a.meta.userAgent()
}
