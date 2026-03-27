package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

var (
	Log *slog.Logger
)

func Init(level string) error {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return err
	}

	cleanOldLogs(logsDir, 7*24*time.Hour)

	logFile := filepath.Join(logsDir, "server-"+time.Now().Format("2006-01-02")+".log")

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	opts := &slog.HandlerOptions{
		Level: lvl,
	}

	handler := slog.NewJSONHandler(file, opts)
	Log = slog.New(handler)
	slog.SetDefault(Log)

	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	Log.Info("logger initialized", "log_file", logFile, "level", level)

	return nil
}

func cleanOldLogs(dir string, maxAge time.Duration) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	cutoff := time.Now().Add(-maxAge)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		info, err := f.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(dir, f.Name()))
		}
	}
}

func Debug(msg string, args ...any) {
	Log.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}

func With(args ...any) *slog.Logger {
	return Log.With(args...)
}
