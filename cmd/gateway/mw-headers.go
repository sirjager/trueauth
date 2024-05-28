package gateway

import (
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

const (
	allow = "allow:"
	block = "block:"
)

func customHeaderMatcher(h []string) func(s string) (string, bool) {
	return func(key string) (string, bool) {
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
