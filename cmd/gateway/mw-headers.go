package gateway

import (
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
)

const (
	allow = "allow:"
	block = "block:"
)

func customHeaderMatcher(l zerolog.Logger, h []string) func(s string) (string, bool) {
	l.Debug().Str("middleware", "allowedHeaders").Msg(REGISTER)
	return func(key string) (string, bool) {
		// logr.Log().Str("middleware", "allowedHeaders").Msg(RUNNING)
		for _, v := range h {
			if v == key {
				if strings.HasPrefix(v, allow) {
					return key[len(allow):], true
				}
				if strings.HasPrefix(v, block) {
					return "", false
				}
			}
		}
		return runtime.DefaultHeaderMatcher(key)
	}
}
