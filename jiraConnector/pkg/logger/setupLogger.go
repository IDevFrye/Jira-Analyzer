package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

const (
	envLocal = "local"
	envDev   = "debug"
	envProd  = "prod"
)

func createLogFile(logFile string) (*os.File, error) {
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	logPath := filepath.Join(projectRoot, "log", logFile)

	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, err
	}

	return os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

func SetupLogger(env string, logFile string) *slog.Logger {
	var log *slog.Logger

	logFilePath, err := createLogFile(logFile)
	if err != nil {
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
		log.Error("can't open log file", Err(err))
	}

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(logFilePath, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(logFilePath, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(logFilePath, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: false}),
		)
	}

	return log
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
