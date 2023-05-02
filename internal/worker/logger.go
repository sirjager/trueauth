package worker

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger zerolog.Logger
}

func NewLogger(logger zerolog.Logger) *Logger {
	return &Logger{logger}
}

func (l *Logger) Print(level zerolog.Level, args ...interface{}) {
	l.logger.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (l *Logger) Debug(args ...interface{}) {
	l.Print(zerolog.DebugLevel, args...)
}

// Info logs a message at Info level.
func (l *Logger) Info(args ...interface{}) {
	l.Print(zerolog.InfoLevel, args...)
}

// Warn logs a message at Warning level.
func (l *Logger) Warn(args ...interface{}) {
	l.Print(zerolog.WarnLevel, args...)
}

// Error logs a message at Error level.
func (l *Logger) Error(args ...interface{}) {
	l.Print(zerolog.ErrorLevel, args...)
}

// Fatal logs a message at Fatal level
// and process will exit with status set to 1.
func (l *Logger) Fatal(args ...interface{}) {
	l.Print(zerolog.FatalLevel, args...)
}
