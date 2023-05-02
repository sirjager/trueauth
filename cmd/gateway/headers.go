package gateway

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

func AllowedHeaders(headers []string) func(s string) (string, bool) {
	// This will take custom headers we want to allow
	return func(key string) (string, bool) {
		for _, h := range headers {
			if h == key {
				return key, true
			}
		}
		return runtime.DefaultHeaderMatcher(key)
	}
}

func GetHeaders(ctx context.Context, key string) ([]string, error) {
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
