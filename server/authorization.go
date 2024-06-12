package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/tokens"
	"google.golang.org/grpc/metadata"

	"github.com/sirjager/trueauth/db/db"
)

type AuthorizedUser struct {
	user    db.User
	payload *tokens.Payload
	token   string
}

const (
	TokenTypeAccess  = "0"
	TokenTypeRefresh = "1"
)

// Checks if request is authenticated and authorized, if not returns error
//
// Extracts authorization header, cookies, verifies and returns AuthorizedUser or error
func (s *Server) authorize(ctx context.Context, refresh ...bool) (auth AuthorizedUser, err error) {
	// NOTE: to allow/block headers, edit cmd/gateway/gateway.go -> incomingHeaders
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = fmt.Errorf("missing metadata from context")
		s.Logr.Error().Err(err).Msg("missing metadata from context")
		return auth, err
	}

	cookies := meta.Get("cookie")
	if len(cookies) > 0 {
		for _, c := range strings.Split(cookies[0], "; ") {
			if parts := strings.Split(c, "="); len(parts) == 2 {
				switch parts[0] {
				case "accessToken":
					auth.token = parts[1]
				case "refreshToken":
					if len(refresh) == 1 && refresh[0] {
						auth.token = parts[1]
					}
				}
			}
		}
	}

	//  extracting authorizationHeader
	values := meta.Get("authorization")
	authHeader := strings.Join(values, ",")
	if fields := strings.Fields(authHeader); len(fields) == 2 {
		auth.token = fields[1]
	}

	if len(auth.token) == 0 {
		s.Logr.Error().Err(err).Msg("missing authorization token")
		return auth, fmt.Errorf(errMissingAuthorization)
	}

	// storing payload and token in authorized if valid token
	incoming, err := s.tokens.VerifyToken(auth.token)
	if err != nil {
		if errors.Is(tokens.ErrExpiredToken, err) {
			s.Logr.Error().Err(err).Msg("expired authorization token")
			return auth, fmt.Errorf(errExpiredToken)
		}
		s.Logr.Error().Err(err).Msg("failed to verify authorization token")
		return auth, fmt.Errorf(errInvalidToken)
	}
	auth.payload = incoming

	var stored tokens.Payload
	_key := tokenKey(incoming.Payload.UserID, incoming.Payload.SessionID, TokenTypeAccess)
	if len(refresh) != 0 && refresh[0] {
		_key = tokenKey(incoming.Payload.UserID, incoming.Payload.SessionID, TokenTypeRefresh)
	}
	if err = s.cache.Get(ctx, _key, &stored); err != nil {
		if errors.Is(cache.ErrNoRecord, err) {
			s.Logr.Error().Err(err).Msg("session not found")
			return auth, fmt.Errorf(errInvalidToken)
		}
		s.Logr.Error().Err(err).Msg("failed to get session")
		return auth, err
	}

	// checking if user exits or not, using user id of token
	auth.user, err = s.store.ReadUser(ctx, stored.Payload.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logr.Error().Err(err).Msg("user not found")
			return auth, fmt.Errorf(errInvalidToken)
		}
		s.Logr.Error().Err(err).Msg("failed to get user")
		return auth, err
	}

	// check if user is verified
	if !auth.user.Verified {
		s.Logr.Error().Msg("email not verified")
		return auth, fmt.Errorf(errEmailNotVerified)
	}

	return auth, err
}
