package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/sirjager/trueauth/pkg/tokens"
)

type MetaData struct {
	payload             *tokens.Payload
	UserAgent           string
	ClientIP            string
	authToken           string
	username            string
	password            string
	authorizationHeader string
}

func (s *Server) extractMetadata(ctx context.Context) *MetaData {
	meta := &MetaData{}

	// For gRPC client-ip
	if p, ok := peer.FromContext(ctx); ok {
		meta.ClientIP = p.Addr.String()
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		//  For HTTP user-agent
		if userAgent := md.Get("grpcgateway-user-agent"); len(userAgent) > 0 {
			meta.UserAgent = userAgent[0]
		}
		//  For HTTP client-ip
		if clientIP := md.Get("x-forwarded-for"); len(clientIP) > 0 {
			meta.ClientIP = clientIP[0]
		}
		//  For gRPC user-agent
		if userAgent := md.Get("user-agent"); len(userAgent) > 0 {
			meta.UserAgent = userAgent[0]
		}

		// For HTTP Authorization Header
		if authHeadr := md.Get("authorization"); len(authHeadr) > 0 {
			meta.authorizationHeader = authHeadr[0]
		}

		// if authorization header is preset,
		// extracting identity and password, if any and tokens
		if meta.authorizationHeader != "" {
			authType := strings.Split(meta.authorizationHeader, " ")[0]

			// we will try to extract access/refresh token if any
			if strings.ToLower(authType) == "bearer" {
				values := strings.Split(meta.authorizationHeader, " ")
				if len(values) == 2 {
					meta.authToken = values[1]
				}
			}

			// we will try to extract identity and password if any
			if strings.ToLower(authType) == "basic" {
				values := strings.Split(meta.authorizationHeader, " ")
				if len(values) == 2 {
					decoded, err := base64.StdEncoding.DecodeString(values[1])
					if err == nil {
						creds := strings.Split(string(decoded), ":")
						if len(creds) == 2 {
							meta.username = creds[0]
							meta.password = creds[1]
						}
					}
				}
			}
		}
	}

	if meta.authToken != "" {
		if payload, err := s.tokens.VerifyToken(meta.authToken); err == nil {
			meta.payload = payload
		}
	}

	return meta
}

func (s *Server) sendHeaders(ctx context.Context, headers map[string]string) error {
	return grpc.SendHeader(ctx, metadata.New(headers))
}

func (s *Server) sendCookies(ctx context.Context, cookies []http.Cookie) error {
	headers := map[string]string{}
	for _, cookie := range cookies {
		bytes, err := json.Marshal(cookie)
		if err != nil {
			return err
		}
		name := "set-cookie:" + cookie.Name
		headers[name] = string(bytes)
	}
	return grpc.SetHeader(ctx, metadata.New(headers))
}
