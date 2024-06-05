package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/tokens"
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
		return auth, fmt.Errorf("missing metadata")
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
		return auth, fmt.Errorf(errMissingAuthorization)
	}

	// storing payload and token in authorized if valid token
	incoming, err := s.tokens.VerifyToken(auth.token)
	if err != nil {
		if errors.Is(tokens.ErrExpiredToken, err) {
			return auth, fmt.Errorf(errExpiredToken)
		}
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
			return auth, fmt.Errorf(errInvalidToken)
		}
		return auth, fmt.Errorf(errInvalidToken)
	}

	// checking if user exits or not, using user id of token
	auth.user, err = s.store.ReadUser(ctx, stored.Payload.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth, fmt.Errorf(errInvalidToken)
		}
		return auth, err
	}

	// check if user is verified
	if !auth.user.Verified {
		return auth, fmt.Errorf(errEmailNotVerified)
	}

	return auth, err
}
