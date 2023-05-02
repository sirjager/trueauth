package gateway

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func Logger(logger zerolog.Logger, handler http.Handler) http.Handler {
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
