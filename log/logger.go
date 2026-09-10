package util

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

var logger = newLogger()

func newLogger() *slog.Logger {
	level := slog.LevelInfo

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}),
	)
}

func Debug(msg string, args ...any) {
	logger.Debug(fmt.Sprintf(msg, args...))
}

func Info(msg string, args ...any) {
	logger.Info(fmt.Sprintf(msg, args...))
}

func Warn(msg string, args ...any) {
	logger.Warn(fmt.Sprintf(msg, args...))
}

func Error(msg string, args ...any) {
	logger.Error(fmt.Sprintf(msg, args...))
}
