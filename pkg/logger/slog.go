package logger

import (
	"context"
	"fmt"
	"github.com/Graylog2/go-gelf/gelf"
	"github.com/gin-gonic/gin"
	"github.com/golang-cz/devslog"
	sloggraylog "github.com/samber/slog-graylog"
	slogmulti "github.com/samber/slog-multi"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	BadKeyId           = "id"
	BadKeyUnderScoreId = "_id"
)

const (
	FormatText = "text"
	FormatJSON = "json"

	LevelTrace = slog.Level(6)
	LevelFatal = slog.Level(12)

	LevelTraceLabel string = "TRACE"
	LevelFatalLabel string = "FATAL"
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace: LevelTraceLabel,
	LevelFatal: LevelFatalLabel,
}

type Fields map[string]any

type Config struct {
	LogLevel           string `yaml:"LEVEL" envconfig:"LOGGER_LEVEL"`
	LoggerFormat       string `yaml:"FORMAT" envconfig:"LOGGER_FORMAT"`
	DebugMode          bool   `yaml:"DEBUG_MODE" envconfig:"LOGGER_DEBUG_MODE"`
	GrayLogActive      bool   `yaml:"GRAYLOG_ACTIVE" envconfig:"LOGGER_GRAYLOG_ACTIVE"`
	GrayLogServer      string `yaml:"GRAYLOG_SERVER" envconfig:"LOGGER_GRAYLOG_SERVER"`
	GrayLogStream      string `yaml:"GRAYLOG_STREAM" envconfig:"LOGGER_GRAYLOG_STREAM"`
	GrayLogRelease     string `yaml:"GRAYLOG_RELEASE" envconfig:"LOGGER_GRAYLOG_RELEASE"`
	GrayLogEnvironment string `yaml:"GRAYLOG_ENVIRONMENT,omitempty" envconfig:"LOGGER_GRAYLOG_ENVIRONMENT,omitempty"`
}

// StandardLogger is the main logger implementation that wraps slog.Logger
type StandardLogger struct {
	Logger *slog.Logger // The underlying slog.Logger instance
	cfg    Config       // Configuration for this logger
}

// Option is a function type that can modify a slog.Logger
type Option func(*slog.Logger)

// WithDefaultFields returns an Option that adds default fields to every log message
// These fields will be included in all log entries created by the logger
func WithDefaultFields(fields Fields) Option {
	return func(logger *slog.Logger) {
		keyVals := make([]any, 0)
		for k, v := range fields {
			keyVals = append(keyVals, k, v)
		}
		*logger = *logger.With(keyVals...)
	}
}

// CreateLogger creates and configures a new StandardLogger based on the provided configuration
// It supports multiple output formats and can send logs to both stdout and Graylog
// Additional options can be provided to customize the logger
func CreateLogger(cfg Config, opts ...Option) *StandardLogger {
	level := parseLogLevel(cfg.LogLevel)

	loggerOptions := &slog.HandlerOptions{Level: level, AddSource: false}
	var handler slog.Handler

	switch cfg.LoggerFormat {
	case FormatJSON:
		handler = slog.NewJSONHandler(os.Stdout, loggerOptions)
	case FormatText:
		handler = slog.NewTextHandler(os.Stdout, loggerOptions)
	default:
		handler = slog.NewTextHandler(os.Stdout, loggerOptions)
	}

	if cfg.GrayLogActive {
		gelfWriter, err := gelf.NewWriter(cfg.GrayLogServer)
		if err != nil {
			log.Fatalf("gelf.NewWriter: %s", err)
		}
		graylogHandler := slog.New(sloggraylog.Option{Level: level, Writer: gelfWriter}.NewGraylogHandler()).
			With("stream", cfg.GrayLogStream)
		handler = slogmulti.Fanout(
			handler,
			graylogHandler.Handler(),
		)
	} else if cfg.DebugMode {
		handler = devslog.NewHandler(os.Stdout, &devslog.Options{
			MaxSlicePrintSize: 4,
			SortKeys:          true,
			TimeFormat:        "[04:05]",
			NewLineAfterLog:   true,
			DebugColor:        devslog.Magenta,
			StringerFormatter: true,
		})
	}

	logger := slog.New(handler)

	for _, opt := range opts {
		opt(logger)
	}

	standardLogger := &StandardLogger{logger, cfg}
	return standardLogger
}

// Trace logs a message at trace level
// This is the most verbose logging level, typically used for detailed debugging information
func (l *StandardLogger) Trace(message string, args ...any) {
	l.Logger.Log(context.Background(), LevelTrace, message, args...)
}

// Debug logs a message at debug level
// This level is used for information that is useful for debugging but not needed in normal operation
func (l *StandardLogger) Debug(message string, args ...any) {
	l.Logger.Debug(message, args...)
}

// Info logs a message at info level
// This is the standard log level for general operational information
func (l *StandardLogger) Info(message string, args ...any) {
	l.Logger.Info(message, args...)
}

// Warn logs a message at warn level
// This level indicates potential issues that don't prevent normal operation
func (l *StandardLogger) Warn(message string, args ...any) {
	l.Logger.Warn(message, args...)
}

// Error logs a message at error level
// This level indicates issues that may require attention but don't cause the application to stop
func (l *StandardLogger) Error(message string, args ...any) {
	l.Logger.Error(message, args...)
}

// Fatal logs a message at error level and then panics
// This should be used for critical errors that prevent the application from continuing
func (l *StandardLogger) Fatal(message string, args ...any) {
	l.Logger.Error(message, args...)
	panic(1)
}

// TraceF logs a formatted message at trace level
// It uses fmt.Sprintf to format the message with the provided arguments
func (l *StandardLogger) TraceF(msg string, args ...any) {
	l.Logger.Log(context.Background(), LevelTrace, fmt.Sprintf(msg, args...))
}

// DebugF logs a formatted message at debug level
// It uses fmt.Sprintf to format the message with the provided arguments
func (l *StandardLogger) DebugF(msg string, args ...any) {
	l.Logger.Debug(fmt.Sprintf(msg, args...))
}

// InfoF logs a formatted message at info level
// It uses fmt.Sprintf to format the message with the provided arguments
func (l *StandardLogger) InfoF(msg string, args ...any) {
	l.Logger.Info(fmt.Sprintf(msg, args...))
}

// WarnF logs a formatted message at warn level
// It uses fmt.Sprintf to format the message with the provided arguments
func (l *StandardLogger) WarnF(msg string, args ...any) {
	l.Logger.Warn(fmt.Sprintf(msg, args...))
}

// ErrorF logs a formatted message at error level
// It uses fmt.Sprintf to format the message with the provided arguments
func (l *StandardLogger) ErrorF(msg string, args ...any) {
	l.Logger.Error(fmt.Sprintf(msg, args...))
}

// FatalF logs a formatted message at error level and then panics
// It uses fmt.Sprintf to format the message with the provided arguments
func (l *StandardLogger) FatalF(msg string, args ...any) {
	l.Logger.Error(fmt.Sprintf(msg, args...))
	panic(1)
}

// TraceWithContext logs a message at trace level with the provided context
// The context can contain values that will be included in the log entry
func (l *StandardLogger) TraceWithContext(ctx context.Context, message string) {
	l.Logger.Log(ctx, LevelTrace, message)
}

// DebugWithContext logs a message at debug level with the provided context
// The context can contain values that will be included in the log entry
func (l *StandardLogger) DebugWithContext(ctx context.Context, message string) {
	l.Logger.DebugContext(ctx, message)
}

// InfoWithContext logs a message at info level with the provided context
// The context can contain values that will be included in the log entry
func (l *StandardLogger) InfoWithContext(ctx context.Context, message string) {
	l.Logger.InfoContext(ctx, message)
}

// WarnWithContext logs a message at warn level with the provided context
// The context can contain values that will be included in the log entry
func (l *StandardLogger) WarnWithContext(ctx context.Context, message string) {
	l.Logger.WarnContext(ctx, message)
}

// ErrorWithContext logs a message at error level with the provided context
// The context can contain values that will be included in the log entry
func (l *StandardLogger) ErrorWithContext(ctx context.Context, message string) {
	l.Logger.ErrorContext(ctx, message)
}

// FatalWithContext logs a message at error level with the provided context and then panics
// The context can contain values that will be included in the log entry
func (l *StandardLogger) FatalWithContext(ctx context.Context, message string) {
	l.Logger.ErrorContext(ctx, message)
	panic(1)
}

// WithField returns a new logger with a single field added to the logging context
// This is useful for adding context to log messages, such as request IDs or user IDs
func (l *StandardLogger) WithField(key string, value any) Logger {
	newLogger := l.Logger.With(key, value)
	return &StandardLogger{Logger: newLogger, cfg: l.cfg}
}

// WithFields returns a new logger with multiple fields added to the logging context
// This allows adding structured data to log messages for better filtering and analysis
func (l *StandardLogger) WithFields(fields Fields) Logger {
	keyvals := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		keyvals = append(keyvals, k, v)
	}
	newLogger := l.Logger.With(keyvals...)
	return &StandardLogger{Logger: newLogger, cfg: l.cfg}
}

// cleanUpDetail removes problematic keys from the detail map
// This prevents issues with certain log processors that have reserved field names
func (l *StandardLogger) cleanUpDetail(detail Detail) Detail {
	badKeys := [...]string{
		BadKeyId,
		BadKeyUnderScoreId,
	}

	for _, key := range badKeys {
		if _, ok := detail[key]; ok {
			l.Warn("Do not use id or _id in your log payload!")
			delete(detail, key)
		}
	}
	return detail
}

// LogApiError logs an API error with context information
// It includes the source file and line where the error occurred, the request URL,
// and additional details provided by the caller
func (l *StandardLogger) LogApiError(context *gin.Context, error error, source *ErrorSource, detail Detail) {
	data := ApiLogStruct{File: source.File, Line: source.Line, Url: context.Request.URL.Path, Detail: l.cleanUpDetail(detail)}
	l.WithField("source", data).
		Error(error.Error())
}

// ApiLogStruct contains structured information about an API error
// It is used to provide context for API error logs
type ApiLogStruct struct {
	Url    string // The URL path of the request
	File   string // The source file where the error occurred
	Line   string // The line number where the error occurred
	Detail Detail // Additional details about the error
}

// Detail is a map type used for storing additional structured information in logs
type Detail map[string]interface{}

// GetErrorSource returns information about the caller's source file and line number
// This is used to provide context about where an error occurred
func GetErrorSource() *ErrorSource {
	_, file, line, _ := runtime.Caller(1)
	return &ErrorSource{
		File: file,
		Line: strconv.Itoa(line),
	}
}

// ErrorSource contains information about where an error occurred in the code
type ErrorSource struct {
	File string // The source file path
	Line string // The line number as a string
}

// parseLogLevel converts a string log level to the corresponding slog.Level
// It supports standard log levels plus custom trace and fatal levels
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "trace":
		return LevelTrace
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "fatal":
		return LevelFatal
	default:
		return slog.LevelInfo
	}
}
