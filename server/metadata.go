package server

import (
	"context"
	"encoding/json"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type MetaData struct {
	userAgent string
	clientIP  string
}

func (s *Server) extractMetadata(ctx context.Context) *MetaData {
	meta := &MetaData{}

	// For gRPC client-ip
	if p, ok := peer.FromContext(ctx); ok {
		meta.clientIP = p.Addr.String()
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		//  For HTTP user-agent
		if userAgent := md.Get("grpcgateway-user-agent"); len(userAgent) > 0 {
			meta.userAgent = userAgent[0]
		}
		//  For HTTP client-ip
		if clientIP := md.Get("x-forwarded-for"); len(clientIP) > 0 {
			meta.clientIP = clientIP[0]
		}
		//  For gRPC user-agent
		if userAgent := md.Get("user-agent"); len(userAgent) > 0 {
			meta.userAgent = userAgent[0]
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
