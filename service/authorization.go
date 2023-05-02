package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/sirjager/trueauth/db/sqlc"
	"github.com/sirjager/trueauth/tokens"
	"google.golang.org/grpc/metadata"
)

const authorizationHeader = "authorization"
const authorizationBearer = "bearer"

type AuthorizedUser struct {
	User    sqlc.User
	Session sqlc.Session
	Token   string
	Payload *tokens.Payload
}

// Checks Authorization header. returns token, payload, account and error
func (s *TrueAuthService) authorize(ctx context.Context) (authorized AuthorizedUser, err error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return authorized, fmt.Errorf("missing metadata")
	}
	values := meta.Get(authorizationHeader)
	if len(values) == 0 {
		return authorized, fmt.Errorf("missing authorization header")
	}

	// authheader will look like:    <token-type> <token>
	// example: Bearer firstpart.secondpart.thirdpart
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return authorized, fmt.Errorf("invalid authorization header format")
	}
	authType := fields[0]
	if strings.ToLower(authType) != authorizationBearer {
		return authorized, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	authorized.Token = fields[1]
	authorized.Payload, err = s.tokens.VerifyToken(authorized.Token)
	if err != nil {
		return authorized, fmt.Errorf("invalid access token: %s", err.Error())
	}

	//? We will also check if token is stored or not
	authorized.Session, err = s.store.Read_Session_ByAccessTokenID(ctx, authorized.Payload.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return authorized, fmt.Errorf("invalid access token")
		}
		return authorized, fmt.Errorf("failed to fetch session: %s", err.Error())
	}

	if authorized.Session.AccessToken != authorized.Token {
		return authorized, fmt.Errorf("invalid access token: %s", err.Error())
	}

	if authorized.Session.UserID.String() != authorized.Payload.Payload.UserID.String() {
		return authorized, fmt.Errorf("invalid access token: %s", err.Error())
	}

	authorized.User, err = s.store.Read_User_ByID(ctx, authorized.Session.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authorized, fmt.Errorf("account does not exists: %s", err.Error())
		}
		return authorized, fmt.Errorf("failed to fetch account: %s", err.Error())
	}

	if s.isUnKnownIP(ctx, authorized.User) {
		return authorized, unknownIPError()
	}

	return authorized, err
}
