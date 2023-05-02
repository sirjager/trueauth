package service

import (
	"context"

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
	ClientIp  string
}

func (s *TrueAuthService) extractMetadata(ctx context.Context) *MetaData {
	meta := &MetaData{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		//  For HTTP user-agent
		if userAgent := md.Get(gatewayUserAgentHeaderKey); len(userAgent) > 0 {
			meta.UserAgent = userAgent[0]
		}
		//  For HTTP client-ip
		if clientIp := md.Get(gatewayClientIpHeaderKey); len(clientIp) > 0 {
			meta.ClientIp = clientIp[0]
		}
		//  For gRPC user-agent
		if userAgent := md.Get(grpcUserAgentHeaderKey); len(userAgent) > 0 {
			meta.UserAgent = userAgent[0]
		}
	}

	// For gRPC client-ip
	if p, ok := peer.FromContext(ctx); ok {
		meta.ClientIp = p.Addr.String()
	}

	return meta
}
