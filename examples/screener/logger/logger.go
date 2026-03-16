package logger

import (
	"log/slog"
	"os"
	"screener/config"
)

type CtxKeyType string

const CtxKey CtxKeyType = "logger"

func New() *slog.Logger {
	options := &slog.HandlerOptions{
		Level: parseLogLevel(config.LogLevel()),
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().UTC().Format("2006-01-02T15:04:05.000000Z"))
			}
			return a
		},
	}

	var handler slog.Handler
	if config.LogFormat() == "json" {
		handler = slog.NewJSONHandler(os.Stdout, options)
	} else {
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	return slog.New(handler)
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	case "info":
		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}
