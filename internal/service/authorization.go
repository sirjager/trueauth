package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/sirjager/trueauth/internal/db/sqlc"
	"github.com/sirjager/trueauth/pkg/tokens"
	"google.golang.org/grpc/metadata"
)

const authorizationHeader = "authorization"
const authorizationBearer = "bearer"

type AuthorizedUser struct {
	user           sqlc.User
	access_token   string
	access_payload *tokens.Payload
	stored_session Session
}

// Checks Authorization header. returns access_token, access_payload, stored_session, user and error
func (s *CoreService) authorize(ctx context.Context) (authorized AuthorizedUser, err error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return AuthorizedUser{}, fmt.Errorf("missing metadata")
	}
	values := meta.Get(authorizationHeader)
	if len(values) == 0 {
		return AuthorizedUser{}, fmt.Errorf("missing authorization header")
	}

	// values[0]: will look like:    <token-type> <token>
	// example: Bearer firstpart.secondpart.thirdpart
	fields := strings.Fields(values[0])
	if len(fields) < 2 {
		return AuthorizedUser{}, fmt.Errorf("invalid authorization header format")
	}
	authType := fields[0]
	if strings.ToLower(authType) != authorizationBearer {
		return AuthorizedUser{}, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	// we are storing incoming bearer token in authorized.Token
	authorized.access_token = fields[1]

	// this will validate token and expiration time
	authorized.access_payload, err = s.tokens.VerifyToken(authorized.access_token)
	if err != nil {
		return AuthorizedUser{}, fmt.Errorf("invalid access token: %s", err.Error())
	}

	//? We will also check if token is stored or not
	session_key := accessTokenKey(authorized.access_payload.Data.UserID, authorized.access_payload.Data.SID)
	session_bytes, err := s.redis.Get(ctx, session_key)

	if err != nil && err != redis.Nil {
		return AuthorizedUser{}, fmt.Errorf("failed to fetch session: %s", err.Error())
	}

	if err == redis.Nil || len(session_bytes) == 0 {
		return AuthorizedUser{}, fmt.Errorf("invalid or expired access token")
	}

	if err = json.Unmarshal(session_bytes, &authorized.stored_session); err != nil {
		return AuthorizedUser{}, err
	}

	// now we match if incoming token and stored token are same or not
	if authorized.access_token != authorized.stored_session.Token {
		return AuthorizedUser{}, fmt.Errorf("invalid access token")
	}

	// Match UserID with incoming token payload
	if authorized.access_payload.Data.UserID != authorized.stored_session.UserID {
		return AuthorizedUser{}, fmt.Errorf("invalid access token")
	}

	authorized.user, err = s.store.ReadUserByID(ctx, []byte(authorized.access_payload.Data.UserID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AuthorizedUser{}, fmt.Errorf("user does not exists: %s", err.Error())
		}
		return AuthorizedUser{}, fmt.Errorf("failed to fetch user: %s", err.Error())
	}

	// This ensures that access token can only be used only by the ip address it was assigned to
	if authorized.access_payload.Data.ClientIP != authorized.stored_session.ClientIP {
		return AuthorizedUser{}, fmt.Errorf("this access token was not assigned to your current ip address, generate new access token")
	}

	if s.isUnKnownIP(ctx, authorized.user) {
		return AuthorizedUser{}, unknownIPError()
	}

	return authorized, err
}
