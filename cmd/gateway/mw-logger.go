package gateway

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

const (
	REGISTER  = "registered"
	EXECUTING = "executing"
	RUNNING   = "running"
)

func logger(logr zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &ResponseRecorder{ResponseWriter: w, StatusCode: 200, Body: &bytes.Buffer{}}
		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		event := logr.Info()

		if rec.StatusCode != http.StatusOK {
			var data map[string]interface{}
			if err := json.Unmarshal(rec.Body.Bytes(), &data); err != nil {
				data = map[string]interface{}{}
			}
			event = logr.Error().Interface("error", data["message"])
		}

		event.
			Str("method", r.Method).
			Str("path", r.RequestURI).
			Str("latency", duration.String()).
			Int("code", int(rec.StatusCode)).
			// Str("status", http.StatusText(rec.StatusCode)).
			Msg("|")
	})
}

type ResponseRecorder struct {
	http.ResponseWriter
	Body       *bytes.Buffer
	StatusCode int
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(b []byte) (int, error) {
	rec.Body.Write(b)
	return rec.ResponseWriter.Write(b)
}
