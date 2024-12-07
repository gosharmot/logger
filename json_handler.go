package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
)

var _ slog.Handler = &(JSONHandler{})

var DefaultLogger = slog.New(&JSONHandler{slog.NewJSONHandler(
	os.Stdout,
	&slog.HandlerOptions{
		AddSource:   true,
		Level:       globalLevel,
		ReplaceAttr: replaceSource,
	},
)})

func NewLogger(out io.Writer, addSource bool) *slog.Logger {
	return slog.New(&JSONHandler{slog.NewJSONHandler(
		out,
		&slog.HandlerOptions{
			AddSource:   addSource,
			Level:       globalLevel,
			ReplaceAttr: replaceSource,
		},
	)})
}

func replaceSource(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		if v, ok := a.Value.Any().(*slog.Source); ok {
			a.Value = slog.AnyValue(&slog.Source{File: v.File, Line: v.Line})
		}
	}
	return a
}

type JSONHandler struct {
	slog.Handler
}

func (h *JSONHandler) Enabled(ctx context.Context, l LogLevel) bool {
	return LevelFromContext(ctx) <= l
}

func (h *JSONHandler) Handle(ctx context.Context, record slog.Record) error {
	if skip, ok := getCallerSkip(ctx); ok {
		var pcs [1]uintptr
		runtime.Callers(skip, pcs[:])
		record.PC = pcs[0]
	}
	return h.Handler.Handle(ctx, record)
}

func (h *JSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &JSONHandler{h.Handler.WithAttrs(attrs)}
}

func (h *JSONHandler) WithGroup(name string) slog.Handler {
	return &JSONHandler{h.Handler.WithGroup(name)}
}
