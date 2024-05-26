package server

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/pkg/validator"
)

type Authenticated struct {
	Profile *db.Profile
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
		return authenticated, fmt.Errorf("missing metadata")
	}

	//  extracting authorizationHeader
	values := meta.Get("authorization")
	if len(values) == 0 {
		return authenticated, fmt.Errorf("missing authorization header")
	}

	// checking auth token format
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return authenticated, fmt.Errorf("invalid authorization header format")
	}

	authType := fields[0]
	if strings.ToLower(authType) != "basic" {
		return authenticated, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	// storing payload and token in authorized if valid token
	decoded, err := base64.StdEncoding.DecodeString(fields[1])
	if err != nil {
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
			return authenticated, fmt.Errorf(errInvalidCredentials)
		}
		return authenticated, fmt.Errorf(errFailedToRetrieveUser, err.Error())
	}

	if !user.Verified {
		return authenticated, fmt.Errorf(errEmailNotVerified)
	}

	if verrifyErr := s.hasher.Verify(user.HashSalt, user.HashPass, password); verrifyErr != nil {
		return authenticated, fmt.Errorf(errInvalidCredentials)
	}

	authenticated.Profile, err = user.Profile()
	if err != nil {
		return authenticated, fmt.Errorf("failed to retrieve profile: %s", err.Error())
	}

	return authenticated, err
}
