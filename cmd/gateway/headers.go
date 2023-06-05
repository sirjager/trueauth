package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

var allowedHeaders = []string{
	"Refresh-Token",
}

func AllowedHeaders() func(s string) (string, bool) {
	// This will take custom headers we want to allow
	return func(key string) (string, bool) {
		for _, h := range allowedHeaders {
			if h == key {
				return key, true
			}
		}
		return runtime.DefaultHeaderMatcher(key)
	}
}
