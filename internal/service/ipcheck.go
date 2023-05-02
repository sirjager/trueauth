package service

import (
	"context"

	"github.com/sirjager/trueauth/internal/db/sqlc"
)

func (s *CoreService) isUnKnownIP(ctx context.Context, user sqlc.User) bool {
	meta := s.extractMetadata(ctx)
	isUnknown := true
	for _, r := range user.AllowedIps {
		if r == meta.ClientIp {
			isUnknown = false
			break
		}
	}
	return isUnknown
}
