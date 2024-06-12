package grpc

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Logger(
	logr zerolog.Logger,
) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		start := time.Now()
		res, err = handler(ctx, req)
		duration := time.Since(start)

		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		event := logr.Info()
		if err != nil {
			event = logr.Error().Err(err)
		}

		event.
			Str("path", info.FullMethod).
			Int("code", int(statusCode)).
			// Str("status", statusCode.String()).
			Str("latency", duration.String()).
			Msg("|")

		return res, err
	}
}

func StreamLogger(
	logr zerolog.Logger,
) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, ss)
		duration := time.Since(start)

		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		event := logr.Info()
		if err != nil {
			event = logr.Error().Err(err)
		}

		event.
			Str("method", info.FullMethod).
			Int("code", int(statusCode)).
			// Str("status", statusCode.String()).
			Str("duration", duration.String()).
			Msg("|")

		return nil
	}
}
