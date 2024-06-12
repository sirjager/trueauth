package server

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/sirjager/gopkg/validator"
	"google.golang.org/grpc/metadata"

	"github.com/sirjager/trueauth/db/db"
)

type Authenticated struct {
	profile *db.Profile
	MetaData
}

const (
	errInvalidCredentials = "invalid credentials"
)

// Used for Signing in with username and password
//
// Checks Authorization header for Basic Auth Only, verifies credentials and returns user and error
func (s *Server) authenticate(ctx context.Context) (authenticated Authenticated, err error) {
	// extracting metadata from FromIncomingContext
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = fmt.Errorf("missing metadata from context")
		s.Logr.Error().Err(err).Msg("missing metadata from context")
		return authenticated, err
	}

	//  extracting authorizationHeader
	values := meta.Get("authorization")
	if len(values) == 0 {
		s.Logr.Error().Err(err).Msg("missing authorization header")
		return authenticated, fmt.Errorf("missing authorization header")
	}

	// checking auth token format
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		s.Logr.Error().Err(err).Msg("invalid authorization header")
		return authenticated, fmt.Errorf("invalid authorization header format")
	}

	authType := fields[0]
	if strings.ToLower(authType) != "basic" {
		s.Logr.Error().Err(err).Msg("unsupported authorization type")
		return authenticated, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	// storing payload and token in authorized if valid token
	decoded, err := base64.StdEncoding.DecodeString(fields[1])
	if err != nil {
		s.Logr.Error().Err(err).Msg("failed to decode base64 authorization")
		return authenticated, fmt.Errorf("failed to decode base64: %s", err.Error())
	}
	credentials := strings.Split(string(decoded), ":")
	identity := credentials[0]
	password := credentials[1]

	findBy := "username"
	if validateErr := validator.ValidateEmail(identity); validateErr == nil {
		findBy = "email"
	}

	var user db.User

	switch findBy {
	case "email":
		user, err = s.store.ReadUserByEmail(ctx, identity)
	default:
		user, err = s.store.ReadUserByUsername(ctx, identity)
	}
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			s.Logr.Error().Err(err).Msg("user not found")
			return authenticated, fmt.Errorf(errInvalidCredentials)
		}
		s.Logr.Error().Err(err).Msg("failed to get user")
		return authenticated, fmt.Errorf(errFailedToRetrieveUser, err.Error())
	}

	if !user.Verified {
		s.Logr.Error().Err(err).Msg("email not verified")
		return authenticated, fmt.Errorf(errEmailNotVerified)
	}

	if err = s.hasher.Verify(user.HashSalt, user.HashPass, password); err != nil {
		s.Logr.Error().Err(err).Msg("invalid password")
		return authenticated, fmt.Errorf(errInvalidCredentials)
	}

	authenticated.profile, err = user.Profile()
	if err != nil {
		s.Logr.Error().Err(err).Msg("failed to build profile")
		return authenticated, fmt.Errorf("failed to retrieve profile: %s", err.Error())
	}

	meta_ := s.extractMetadata(ctx)
	authenticated.clientIP = meta_.clientIP
	authenticated.userAgent = meta_.userAgent

	return authenticated, err
}
