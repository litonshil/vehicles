package logger

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func InitLogger() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	logger = slog.New(logHandler)
	slog.SetDefault(logger)
}
