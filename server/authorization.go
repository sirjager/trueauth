package server

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/pkg/tokens"
	rpc "github.com/sirjager/trueauth/stubs"
)

type AuthorizedUser struct {
	User    db.User
	Session db.Session
	Payload *tokens.Payload
	Profile *rpc.User
	Meta    *MetaData
	Token   string
}

// Checks if request is authenticated and authorized, if not returns error
// 
// Extracts authorization header, cookies, verifies and returns AuthorizedUser or error
func (s *Server) authorize(ctx context.Context, refresh ...bool) (auth AuthorizedUser, err error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return auth, fmt.Errorf("missing metadata")
	}

	cookies := meta.Get("grpcgateway-cookie")
	if len(cookies) > 0 {
		for _, c := range strings.Split(cookies[0], "; ") {
			if parts := strings.Split(c, "="); len(parts) == 2 {
				switch parts[0] {
				case "accessToken":
					auth.Token = parts[1]
				case "refreshToken":
					if len(refresh) == 1 && refresh[0] {
						auth.Token = parts[1]
					}
				}
			}
		}
	}

	//  extracting authorizationHeader
	values := meta.Get("authorization")
	if len(values) != 0 {
		authHeader := values[0]
		if fields := strings.Fields(authHeader); len(fields) == 2 {
			auth.Token = fields[1]
		}
	}

	if len(auth.Token) == 0 {
		return auth, fmt.Errorf(errMissingAuthorization)
	}

	// storing payload and token in authorized if valid token
	auth.Payload, err = s.tokens.VerifyToken(auth.Token)
	if err != nil {
		return auth, fmt.Errorf(errUnauthorized)
	}

	switch auth.Payload.Payload.Type {
	case "access":
		auth.Session, err = s.store.ReadSessionByAccessTokenID(ctx, auth.Payload.ID)
	case "refresh":
		auth.Session, err = s.store.ReadSession(ctx, auth.Payload.ID)
	}

	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return auth, fmt.Errorf(errUnauthorized)
		}
		return auth, fmt.Errorf("failed to fetch session: %s", err.Error())
	}

	switch auth.Payload.Payload.Type {
	case "access":
		if auth.Session.AccessToken != auth.Token {
			return auth, fmt.Errorf(errUnauthorized)
		}
	case "refresh":
		if auth.Session.RefreshToken != auth.Token {
			return auth, fmt.Errorf(errUnauthorized)
		}
	}

	// checking user id of stored tokens and incoming token is same or not
	if !bytes.Equal(auth.Session.UserID, auth.Payload.Payload.UserID) {
		return auth, fmt.Errorf(errUnauthorized)
	}

	// checking if user exits or not, using user id of token
	auth.User, err = s.store.ReadUser(ctx, auth.Session.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth, fmt.Errorf(errUnauthorized)
		}
		return auth, fmt.Errorf(errFailedToRetrieveUser, err.Error())
	}

	// check if user is verified
	if !auth.User.Verified {
		return auth, fmt.Errorf(errEmailNotVerified)
	}

	profile, err := auth.User.Profile()
	if err != nil {
		return auth, fmt.Errorf("something went wrong")
	}
	auth.Profile = publicProfile(profile)

	auth.Meta = s.extractMetadata(ctx)

	// finally return authorized and error if any
	return auth, err
}
