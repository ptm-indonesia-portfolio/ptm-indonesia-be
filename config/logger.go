package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

func NewLogger(cfg *AppConfig) (*logrus.Logger, func(), error) {
	if err := os.MkdirAll(filepath.Dir(cfg.Log.FilePath), 0o755); err != nil {
		return nil, nil, fmt.Errorf("create log directory: %w", err)
	}

	logFile, err := os.OpenFile(cfg.Log.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file: %w", err)
	}

	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		_ = logFile.Close()
		return nil, nil, fmt.Errorf("parse log level: %w", err)
	}

	if level > logrus.ErrorLevel {
		level = logrus.ErrorLevel
	}

	logger := logrus.New()
	logger.SetOutput(logFile)
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})

	cleanup := func() {
		_ = logFile.Close()
	}

	return logger, cleanup, nil
}
