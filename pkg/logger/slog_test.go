package logger

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
)

func TestCreateLogger(t *testing.T) {
	cfg := Config{
		LogLevel:     "info",
		LoggerFormat: FormatText,
		DebugMode:    false,
	}

	logger := CreateLogger(cfg)
	if logger == nil {
		t.Fatal("Expected logger instance, got nil")
	}
}

func TestLoggerLevels(t *testing.T) {
	cfg := Config{LogLevel: "debug", LoggerFormat: FormatText}
	logger := CreateLogger(cfg)

	logger.Trace("trace message")
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
}

func TestLoggerWithFields(t *testing.T) {
	cfg := Config{LogLevel: "info", LoggerFormat: FormatText}
	logger := CreateLogger(cfg)

	newLogger := logger.WithFields(Fields{"key": "value"})
	if newLogger == nil {
		t.Fatal("Expected logger instance with fields, got nil")
	}
}

func TestLoggerWithContext(t *testing.T) {
	cfg := Config{LogLevel: "info", LoggerFormat: FormatText}
	logger := CreateLogger(cfg)
	ctx := context.Background()

	logger.InfoWithContext(ctx, "info with context")
}

func TestLoggerWithDefaultFields(t *testing.T) {
	cfg := Config{LogLevel: "info", LoggerFormat: FormatText}
	logger := CreateLogger(cfg, WithDefaultFields(Fields{"service": "test"}))

	logger.Info("Test log with default fields")
}

func TestParseLogLevel(t *testing.T) {
	tests := map[string]slog.Level{
		"trace": LevelTrace,
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
		"fatal": LevelFatal,
	}

	for input, expected := range tests {
		result := parseLogLevel(input)
		if result != expected {
			t.Errorf("Expected %v, got %v for input %s", expected, result, input)
		}
	}
}

func TestLoggerOutputText(t *testing.T) {
	var buf bytes.Buffer
	cfg := Config{LogLevel: "info", LoggerFormat: FormatText}
	logger := CreateLogger(cfg)
	logger.Logger = slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Info("test message text format")

	if !bytes.Contains(buf.Bytes(), []byte("test message text format")) {
		t.Errorf("Expected log message not found in text output: %s", buf.String())
	}
}

func TestLoggerOutputJSON(t *testing.T) {
	var buf bytes.Buffer
	cfg := Config{LogLevel: "info", LoggerFormat: FormatJSON}
	logger := CreateLogger(cfg)
	logger.Logger = slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Info("test message json format")

	if !bytes.Contains(buf.Bytes(), []byte("test message json format")) {
		t.Errorf("Expected log message not found in JSON output: %s", buf.String())
	}
}
