package service

import (
	"context"

	"github.com/sirjager/trueauth/db/sqlc"
)

func (s *TrueAuthService) isUnKnownIP(ctx context.Context, account sqlc.User) bool {
	meta := s.extractMetadata(ctx)
	isUnknown := true
	for _, r := range account.AllowedIps {
		if r == meta.ClientIp {
			isUnknown = false
			break
		}
	}
	return isUnknown
}
