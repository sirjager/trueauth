package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
)

const (
	Reset      = "\033[0m"
	Black      = "\033[0;30m"
	Red        = "\033[0;31m"
	Green      = "\033[0;32m"
	Yellow     = "\033[0;33m"
	Blue       = "\033[0;34m"
	Purple     = "\033[0;35m"
	Cyan       = "\033[0;36m"
	White      = "\033[0;37m"
	Gray       = "\033[0;90m"
	BoldBlack  = "\033[1;30m"
	BoldRed    = "\033[1;31m"
	BoldGreen  = "\033[1;32m"
	BoldYellow = "\033[1;33m"
	BoldBlue   = "\033[1;34m"
	BoldPurple = "\033[1;35m"
	BoldCyan   = "\033[1;36m"
	BoldWhite  = "\033[1;37m"
)

type Config struct {
	ServerName string `mapstructure:"SERVER_NAME"`
	Logfile    string `mapstructure:"LOGFILE"`
}

type Logger struct {
	file *os.File
	Logr zerolog.Logger
}

func NewLogger(cfg Config) (*Logger, error) {
	var logr zerolog.Logger
	var file *os.File
	var err error

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	if cfg.Logfile != "" {
		file, err = os.OpenFile(cfg.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		logr = zerolog.New(file)
	} else {
		logr = zerolog.New(newConsoleWriter())
	}

	logr = logr.With().Str("server", cfg.ServerName).Logger()
	logr = logr.With().Timestamp().Caller().Logger()

	return &Logger{Logr: logr, file: file}, nil
}

func (l *Logger) Close() error {
	if l.file != nil {
		l.Logr.Info().Msg("stopping the logger and closing logfile")
		return l.file.Close()
	}
	return nil
}

func newConsoleWriter() io.Writer {
	return zerolog.ConsoleWriter{
		Out:         os.Stderr,
		NoColor:     false,
		FieldsOrder: []string{"server", "method", "code", "status", "latency", "path"},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("%s%s=%s", Purple, i, Reset)
		},
		FormatCaller: func(i interface{}) string {
			return fmt.Sprintf("%s%s%s", Blue, i, Reset)
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%s%s%s", Gray, i, Reset)
		},
		FormatErrFieldName: func(i interface{}) string {
			return fmt.Sprintf("%s%s=%s", Red, i, Reset)
		},
		FormatErrFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%s%s%s", Red, i, Reset)
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%s%s%s", White, i, Reset)
		},
		FormatTimestamp: func(i interface{}) string {
			return fmt.Sprintf("%s%s%s", Gray, i, Reset)
		},
	}
}
