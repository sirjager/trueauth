package service

import (
	"context"

	"github.com/sirjager/trueauth/db/sqlc"
)

func (s *TrueAuthService) isBlockedIP(iprecord sqlc.Iprecord, ctx context.Context) bool {
	meta := s.extractMetadata(ctx)
	isBlockedIP := false
	for _, r := range iprecord.BlockedIps {
		if r == meta.ClientIp {
			isBlockedIP = true
			break
		}
	}
	return isBlockedIP
}

func (s *TrueAuthService) isKnownIP(iprecord sqlc.Iprecord, ctx context.Context) bool {
	meta := s.extractMetadata(ctx)
	isNewIP := false
	for _, r := range iprecord.AllowedIps {
		if r == meta.ClientIp {
			isNewIP = true
			break
		}
	}
	return isNewIP
}
