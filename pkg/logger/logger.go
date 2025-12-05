package logger

import "context"

// Logger is the interface that defines all logging methods
// It provides methods for different log levels, formatted logging, context-aware logging,
// and methods to add fields to the logging context
type Logger interface {
	Trace(msg string, args ...any)
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)

	TraceF(msg string, args ...any)
	DebugF(msg string, args ...any)
	InfoF(msg string, args ...any)
	WarnF(msg string, args ...any)
	ErrorF(msg string, args ...any)
	FatalF(msg string, args ...any)

	TraceWithContext(ctx context.Context, msg string)
	DebugWithContext(ctx context.Context, msg string)
	InfoWithContext(ctx context.Context, msg string)
	WarnWithContext(ctx context.Context, msg string)
	ErrorWithContext(ctx context.Context, msg string)
	FatalWithContext(ctx context.Context, msg string)

	WithField(key string, value any) Logger
	WithFields(fields Fields) Logger
}
