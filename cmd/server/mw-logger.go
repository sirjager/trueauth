package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func HTTPLogger(logger zerolog.Logger, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}

		handler.ServeHTTP(rec, req)
		duration := time.Since(start)

		event := logger.Info()
		if rec.StatusCode != http.StatusOK {
			var data map[string]interface{}
			if err := json.Unmarshal(rec.Body, &data); err != nil {
				data = map[string]interface{}{}
			}
			event = logger.Error().Interface("error", data["message"])
		}

		event.
			Str("protocol", "REST").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Str("duration", duration.String()).
			Int("code", int(rec.StatusCode)).
			Str("status", http.StatusText(rec.StatusCode)).
			Msg("")
	})
}

func GRPCLogger(logger zerolog.Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		start := time.Now()
		res, err = handler(ctx, req)
		duration := time.Since(start)

		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		event := logger.Info()
		if err != nil {
			event = logger.Error().Err(err)
		}

		event.
			Str("protocol", "gRPC").
			Str("method", info.FullMethod).
			Int("code", int(statusCode)).
			Str("status", statusCode.String()).
			Str("duration", duration.String()).
			Msg("")

		return res, err

	}
}

func GRPCStreamLogger(logger zerolog.Logger) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, ss)
		duration := time.Since(start)

		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		event := logger.Info()
		if err != nil {
			event = logger.Error().Err(err)
		}

		event.
			Str("protocol", "gRPC").
			Str("method", info.FullMethod).
			Int("code", int(statusCode)).
			Str("status", statusCode.String()).
			Str("duration", duration.String()).
			Msg("")

		return nil
	}
}
