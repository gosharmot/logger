package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync/atomic"
)

var (
	global atomic.Pointer[slog.Logger]
)

// need for testing
var exitF = os.Exit

func init() {
	SetLevel(LevelInfo)
	SetLogger(DefaultLogger)
}

func SetLogger(l *slog.Logger) {
	global.Store(l)
}

func SetLevel(l LogLevel) {
	globalLevel.SetLevel(l)
}

func Level() LogLevel {
	return globalLevel.Level()
}

func Debug(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).DebugContext(setCallerSkip(ctx, 5), msg, args...)
}

func Debugf(ctx context.Context, format string, args ...any) {
	FromContext(ctx).DebugContext(setCallerSkip(ctx, 5), fmt.Sprintf(format, args...))
}

func Info(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).InfoContext(setCallerSkip(ctx, 5), msg, args...)
}

func Infof(ctx context.Context, format string, args ...any) {
	FromContext(ctx).InfoContext(setCallerSkip(ctx, 5), fmt.Sprintf(format, args...))
}

func Warn(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).WarnContext(setCallerSkip(ctx, 5), msg, args...)
}

func Warnf(ctx context.Context, format string, args ...any) {
	FromContext(ctx).WarnContext(setCallerSkip(ctx, 5), fmt.Sprintf(format, args...))
}

func Error(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).ErrorContext(setCallerSkip(ctx, 5), msg, args...)
}

func Errorf(ctx context.Context, format string, args ...any) {
	FromContext(ctx).ErrorContext(setCallerSkip(ctx, 5), fmt.Sprintf(format, args...))
}

func Fatal(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).ErrorContext(setCallerSkip(ctx, 5), msg, args...)
	exitF(1)
}

func Fatalf(ctx context.Context, format string, args ...any) {
	FromContext(ctx).ErrorContext(setCallerSkip(ctx, 5), fmt.Sprintf(format, args...))
	exitF(1)
}

func WithGroup(ctx context.Context, group string) *slog.Logger {
	return FromContext(ctx).WithGroup(group)
}

func With(ctx context.Context, args ...any) context.Context {
	return ToContext(ctx, FromContext(ctx).With(args...))
}
