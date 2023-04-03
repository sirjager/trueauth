package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/sirjager/trueauth/db/sqlc"
	"google.golang.org/grpc/metadata"
)

const authorizationHeader = "authorization"
const authorizationBearer = "bearer"

// Checks Authorization header. returns token, payload, user and error
func (s *TrueAuthService) authorize(ctx context.Context) (user sqlc.User, ipRecord sqlc.Iprecord, err error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return sqlc.User{}, sqlc.Iprecord{}, fmt.Errorf("missing metadata")
	}
	values := meta.Get(authorizationHeader)
	if len(values) == 0 {
		return sqlc.User{}, sqlc.Iprecord{}, fmt.Errorf("missing authorization header")
	}

	// authheader will look like:    <token-type> <token>
	// example: Bearer firstpart.secondpart.thirdpart
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return sqlc.User{}, sqlc.Iprecord{}, fmt.Errorf("invalid authorization header format")
	}
	authType := fields[0]
	if strings.ToLower(authType) != authorizationBearer {
		return sqlc.User{}, sqlc.Iprecord{}, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	tokenString := fields[1]
	tokenPayload, err := s.tokens.VerifyToken(tokenString)
	if err != nil {
		return sqlc.User{}, sqlc.Iprecord{}, fmt.Errorf("invalid access token: %s", err.Error())
	}

	//? We will also check if token is stored or not
	session, err := s.store.GetSessionByAccessTokenID(ctx, tokenPayload.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return sqlc.User{}, sqlc.Iprecord{}, fmt.Errorf("invalid access token: %s", err.Error())
		}
		return sqlc.User{}, sqlc.Iprecord{}, fmt.Errorf("failed to fetch session: %s", err.Error())
	}

	if session.AccessToken != tokenString {
		return sqlc.User{}, sqlc.Iprecord{}, fmt.Errorf("invalid access token: %s", err.Error())
	}

	if session.UserID.String() != tokenPayload.Payload.UserID.String() {
		return sqlc.User{}, sqlc.Iprecord{}, fmt.Errorf("invalid access token: %s", err.Error())
	}

	//
	ipRecord, err = s.store.GetIPRecordByUserID(ctx, session.UserID)
	if err != nil {
		return sqlc.User{}, sqlc.Iprecord{}, fmt.Errorf("failed to fetch ip records: %s", err.Error())
	}

	if s.isBlockedIP(ipRecord, ctx) {
		return sqlc.User{}, sqlc.Iprecord{}, fmt.Errorf("ip address is blocked")
	}

	user, err = s.store.GetUserByID(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, sqlc.Iprecord{}, fmt.Errorf("user does not exists: %s", err.Error())
		}
		return sqlc.User{}, sqlc.Iprecord{}, fmt.Errorf("failed to fetch user: %s", err.Error())
	}

	return user, ipRecord, err
}
