package logger

import (
	"log/slog"
	"os"
)

func SetupLogger() *slog.Logger {

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	logger := slog.New(handler)
    
    slog.SetDefault(logger)
    
    return logger
}
