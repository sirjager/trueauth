package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/sirjager/gopkg/cache"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/internal/tokens"
	rpc "github.com/sirjager/trueauth/rpc"
)

type AuthorizedUser interface {
	Profile() *rpc.User
	ClientIP() string
	UserAgent() string
	Token() string
	Payload() tokens.Payload
	User () db.User
}

type _authorizedUser struct {
	user    *rpc.User
	dbuser  db.User
	payload tokens.Payload
	token   string
	meta    MetaData
}

const (
	TokenTypeAccess  = "0"
	TokenTypeRefresh = "1"
)

// Checks if request is authenticated and authorized, if not returns error
//
// Extracts authorization header, cookies, verifies and returns AuthorizedUser or error
func (s *Server) authorize(ctx context.Context, refresh ...bool) (AuthorizedUser, error) {
	var err error
	// NOTE: to allow/block headers, edit cmd/gateway/gateway.go -> incomingHeaders
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = fmt.Errorf("missing metadata from context")
		s.Logr.Error().Err(err).Msg("missing metadata from context")
		return nil, err
	}

	auth := &_authorizedUser{}

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
		// return error if no authorization header
		s.Logr.Error().Err(err).Msg("missing authorization token")
		return nil, fmt.Errorf(errMissingAuthorization)
	}

	// storing payload and token in authorized if valid token
	incoming, err := s.tokens.VerifyToken(auth.token)
	if err != nil {
		if errors.Is(tokens.ErrExpiredToken, err) {
			// return error if expired token
			s.Logr.Error().Err(err).Msg("expired authorization token")
			return nil, fmt.Errorf(errExpiredToken)
		}
		// return error if invalid token
		s.Logr.Error().Err(err).Msg("failed to verify authorization token")
		return nil, fmt.Errorf(errInvalidToken)
	}
	auth.payload = *incoming

	var stored tokens.Payload
	_key := tokenKey(incoming.Payload.UserID, incoming.Payload.SessionID, TokenTypeAccess)
	if len(refresh) != 0 && refresh[0] {
		_key = tokenKey(incoming.Payload.UserID, incoming.Payload.SessionID, TokenTypeRefresh)
	}
	if err = s.cache.Get(ctx, _key, &stored); err != nil {
		if errors.Is(cache.ErrNoRecord, err) {
			// return error if session not found
			s.Logr.Error().Err(err).Msg("session not found")
			return nil, fmt.Errorf(errInvalidToken)
		}
		s.Logr.Error().Err(err).Msg("failed to get session")
		return nil, err
	}

	// checking if user exits or not, using user id of token
	auth.dbuser, err = s.store.ReadUser(ctx, stored.Payload.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// return error invalid token even if user not found
			s.Logr.Error().Err(err).Msg("user not found")
			return nil, fmt.Errorf(errInvalidToken)
		}
		// return error if failed to get user from db
		s.Logr.Error().Err(err).Msg("failed to get user")
		return nil, err
	}

	// check if user is verified
	if !auth.dbuser.Verified {
		// return error if email not verified
		s.Logr.Error().Msg("email not verified")
		return nil, fmt.Errorf(errEmailNotVerified)
	}

	auth.meta = s.extractMetadata(ctx)
	auth.user = &rpc.User{
		Id:        auth.dbuser.ID,
		Email:     auth.dbuser.Email,
		Username:  auth.dbuser.Username,
		Firstname: auth.dbuser.Firstname,
		Lastname:  auth.dbuser.Lastname,
		Verified:  auth.dbuser.Verified,
		CreatedAt: timestamppb.New(auth.dbuser.CreatedAt),
		UpdatedAt: timestamppb.New(auth.dbuser.UpdatedAt),
	}

	return auth, err
}

func (a *_authorizedUser) User() db.User {
	return a.dbuser
}

func (a *_authorizedUser) Profile() *rpc.User {
	return a.user
}

func (a *_authorizedUser) ClientIP() string {
	return a.meta.clientIP()
}

func (a *_authorizedUser) UserAgent() string {
	return a.meta.userAgent()
}

func (a *_authorizedUser) Token() string {
	return a.token
}

func (a *_authorizedUser) Payload() tokens.Payload {
	return a.payload
}
