package logger

import "log/slog"

// Error adds an error attribute to the log entry.
func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}

// Operation adds an operation attribute to the log entry.
func Operation(operation string) slog.Attr {
	return slog.String("operation", operation)
}
