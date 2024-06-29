package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type MetaData interface {
	clientIP() string
	userAgent() string
}

type _MetaData struct {
	_userAgent string
	_clientIP  string
}

func (s *Server) extractMetadata(ctx context.Context) MetaData {
	meta := &_MetaData{}

	// For gRPC client-ip
	if p, ok := peer.FromContext(ctx); ok {
		meta._clientIP = p.Addr.String()
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		//  For HTTP user-agent
		if userAgent := md.Get("grpcgateway-user-agent"); len(userAgent) > 0 {
			meta._userAgent = userAgent[0]
		}
		//  For HTTP client-ip
		if clientIP := md.Get("x-forwarded-for"); len(clientIP) > 0 {
			meta._clientIP = clientIP[0]
		}
		//  For gRPC user-agent
		if userAgent := md.Get("user-agent"); len(userAgent) > 0 {
			meta._userAgent = userAgent[0]
		}
	}

	return meta
}

func (m *_MetaData) clientIP() string {
	return m._clientIP
}

func (m *_MetaData) userAgent() string {
	return m._userAgent
}

func (s *Server) SetStatusCode(ctx context.Context, code int) error {
	return grpc.SendHeader(ctx, metadata.Pairs("x-http-code", fmt.Sprintf("%d", code)))
}

func (s *Server) SetHeaders(ctx context.Context, headers map[string]string) error {
	return grpc.SendHeader(ctx, metadata.New(headers))
}

func (s *Server) SetCookies(ctx context.Context, cookies []http.Cookie) error {
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
