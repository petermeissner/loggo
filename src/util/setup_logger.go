package util

import (
	"log/slog"
	"os"
)

func Setup_logger(log_level string, log_type string) *slog.Logger {
	// standard options for logger
	logger_opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	}

	// modify options according to log_level
	switch log_level {
	case "debug":
		logger_opts.Level = slog.LevelDebug
		logger_opts.AddSource = true
	case "info":
		logger_opts.Level = slog.LevelInfo
	case "warn":
		logger_opts.Level = slog.LevelWarn
	default:
		logger_opts.Level = slog.LevelError
	}

	// modify options according to log_type adn return logger
	if log_type == "json" {
		return slog.New(slog.NewJSONHandler(os.Stdout, logger_opts))
	} else {
		return slog.New(slog.NewTextHandler(os.Stdout, logger_opts))
	}
}
