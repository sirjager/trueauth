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

func logger(logr zerolog.Logger, handler http.Handler) http.Handler {
	logr.Info().Str("middleware", "logger").Msg(REGISTER)
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		rec := &ResponseRecorder{ResponseWriter: res, StatusCode: 200, Body: &bytes.Buffer{}}
		handler.ServeHTTP(rec, req)
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
			Str("protocol", "REST").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Str("duration", duration.String()).
			Int("code", int(rec.StatusCode)).
			Str("status", http.StatusText(rec.StatusCode)).
			Msg("")
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
