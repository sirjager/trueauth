package grpc

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	gatewayUserAgentHeaderKey = "grpcgateway-user-agent"
	gatewayClientIPHeaderKey  = "x-forwarded-for"
	grpcUserAgentHeaderKey    = "user-agent"
	authorizationHeader       = "authorization"
	authorizationBearer       = "bearer"
)

type MetaData struct {
	// Payload is the token payload of AuthToken
	UserAgent          string
	ClientIP           string
	AuthToken          string
	Username           string
	Password           string
	autorizationHeader string
}

const (
	ErrInvalidAuthorizationHeader = "invalid authorization header"
	ErrDecodingBasicAuth          = "error decoding basic auth"
)

func ExtractMetadata(ctx context.Context) (meta *MetaData, err error) {
	meta = &MetaData{}

	// NOTE: Extract Grpc Client IP
	if p, ok := peer.FromContext(ctx); ok {
		meta.ClientIP = p.Addr.String()
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {

		// NOTE: Extaact HTTP UserAgent
		if userAgent := md.Get(gatewayUserAgentHeaderKey); len(userAgent) > 0 {
			meta.UserAgent = userAgent[0]
		}

		// NOTE: Extract HTTP Client IP
		if clientIP := md.Get(gatewayClientIPHeaderKey); len(clientIP) > 0 {
			meta.ClientIP = clientIP[0]
		}

		// NOTE: Extract Grpc UserAgent
		if userAgent := md.Get(grpcUserAgentHeaderKey); len(userAgent) > 0 {
			meta.UserAgent = userAgent[0]
		}

		// NOTE: Extract Authorization Header if any
		if authHeadr := md.Get(authorizationHeader); len(authHeadr) > 0 {
			meta.autorizationHeader = authHeadr[0]
		}

		// NOTE: Extract Authorization Brearer / Basic Auth if present
		//
		// If Authorization header is present then,
		// we will try to extract bearer token, basic auth if present
		// we will also try to valiate and return proper errors
		if meta.autorizationHeader != "" {
			authType := strings.Split(meta.autorizationHeader, " ")[0]
			// we will try to extract bearer token if any
			if strings.ToLower(authType) == "bearer" {
				values := strings.Split(meta.autorizationHeader, " ")
				if len(values) != 2 {
					return meta, fmt.Errorf(ErrInvalidAuthorizationHeader)
				}
				meta.AuthToken = values[1]
			}

			// we will try to extract username and password if basic auth
			if strings.ToLower(authType) == "basic" {
				values := strings.Split(meta.autorizationHeader, " ")
				if len(values) != 2 {
					return meta, fmt.Errorf(ErrInvalidAuthorizationHeader)
				}
				decoded, err := base64.StdEncoding.DecodeString(values[1])
				if err != nil {
					return meta, fmt.Errorf(ErrDecodingBasicAuth)
				}
				creds := strings.Split(string(decoded), ":")
				if len(creds) != 2 {
					return meta, fmt.Errorf(ErrInvalidAuthorizationHeader)
				}
				meta.Username = creds[0]
				meta.Password = creds[1]
			}

			//
		}
	}
	return
}
