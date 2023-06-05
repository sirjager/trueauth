package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	gatewayUserAgentHeaderKey = "grpcgateway-user-agent"
	gatewayClientIpHeaderKey  = "x-forwarded-for"

	grpcUserAgentHeaderKey = "user-agent"
)

type MetaData struct {
	UserAgent string
	ClientIP  string
}

func (s *CoreService) extractMetadata(ctx context.Context) *MetaData {
	meta := &MetaData{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		//  For HTTP user-agent
		if userAgent := md.Get(gatewayUserAgentHeaderKey); len(userAgent) > 0 {
			meta.UserAgent = userAgent[0]
		}
		//  For HTTP client-ip
		if clientIp := md.Get(gatewayClientIpHeaderKey); len(clientIp) > 0 {
			meta.ClientIP = clientIp[0]
		}
		//  For gRPC user-agent
		if userAgent := md.Get(grpcUserAgentHeaderKey); len(userAgent) > 0 {
			meta.UserAgent = userAgent[0]
		}
	}

	// For gRPC client-ip
	if p, ok := peer.FromContext(ctx); ok {
		meta.ClientIP = p.Addr.String()
	}

	return meta
}

func (s *CoreService) extractHeaders(ctx context.Context, key string) ([]string, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	values := meta.Get(key)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing %s from headers", key)
	}

	return values, nil
}
