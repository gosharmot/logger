package logger

import (
	"context"
	"log/slog"
)

type contextKey int

const (
	loggerKey contextKey = iota
	levelKey
	callerScipKey
)

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	l := global.Load()
	return l
}

func ToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func LevelFromContext(ctx context.Context) LogLevel {
	if level, ok := ctx.Value(levelKey).(LogLevel); ok {
		return level
	}
	return globalLevel.Level()
}

func LevelToContext(ctx context.Context, l LogLevel) context.Context {
	return context.WithValue(ctx, levelKey, l)
}

func setCallerSkip(ctx context.Context, skip int) context.Context {
	return context.WithValue(ctx, callerScipKey, skip)
}

func getCallerSkip(ctx context.Context) (int, bool) {
	skip, ok := ctx.Value(callerScipKey).(int)
	return skip, ok
}
