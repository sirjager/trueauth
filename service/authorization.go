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

// Checks Authorization header. returns token, payload, account and error
func (s *TrueAuthService) authorize(ctx context.Context) (account sqlc.Account, ip sqlc.Ip, err error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return account, ip, fmt.Errorf("missing metadata")
	}
	values := meta.Get(authorizationHeader)
	if len(values) == 0 {
		return account, ip, fmt.Errorf("missing authorization header")
	}

	// authheader will look like:    <token-type> <token>
	// example: Bearer firstpart.secondpart.thirdpart
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return account, ip, fmt.Errorf("invalid authorization header format")
	}
	authType := fields[0]
	if strings.ToLower(authType) != authorizationBearer {
		return account, ip, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	tokenString := fields[1]
	tokenPayload, err := s.tokens.VerifyToken(tokenString)
	if err != nil {
		return account, ip, fmt.Errorf("invalid access token: %s", err.Error())
	}

	//? We will also check if token is stored or not
	session, err := s.store.GetSessionByAccessTokenID(ctx, tokenPayload.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return account, ip, fmt.Errorf("invalid access token: %s", err.Error())
		}
		return account, ip, fmt.Errorf("failed to fetch session: %s", err.Error())
	}

	if session.AccessToken != tokenString {
		return account, ip, fmt.Errorf("invalid access token: %s", err.Error())
	}

	if session.AccountID.String() != tokenPayload.Payload.AccountID.String() {
		return account, ip, fmt.Errorf("invalid access token: %s", err.Error())
	}

	//
	ip, err = s.store.GetIPByAccountID(ctx, session.AccountID)
	if err != nil {
		return account, sqlc.Ip{}, fmt.Errorf("failed to fetch ip records: %s", err.Error())
	}

	if s.isBlockedIP(ip, ctx) {
		return account, sqlc.Ip{}, fmt.Errorf("ip address is blocked")
	}

	account, err = s.store.GetAccountByID(ctx, session.AccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Account{}, sqlc.Ip{}, fmt.Errorf("account does not exists: %s", err.Error())
		}
		return sqlc.Account{}, sqlc.Ip{}, fmt.Errorf("failed to fetch account: %s", err.Error())
	}

	return account, ip, err
}
