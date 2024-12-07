package logger

import (
	"bytes"
	"context"
	"log/slog"
	"reflect"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("default_from_context", func(t *testing.T) {
		if FromContext(context.Background()) != DefaultLogger {
			t.Fatal("logger form context not equal default logger")
		}
	})

	t.Run("set_default", func(t *testing.T) {
		defer SetLogger(DefaultLogger)

		l := NewLogger(nil, false)
		SetLogger(l)
		if FromContext(context.Background()) != l {
			t.Fatal("logger form context not equal established logger")
		}
	})

	t.Run("set_to_context", func(t *testing.T) {
		l := NewLogger(nil, false)
		ctx := ToContext(context.Background(), l)
		if FromContext(ctx) != l {
			t.Fatal("logger form context not equal established logger")
		}
		if global.Load() == l {
			t.Fatal("global logger must be default")
		}
	})

	t.Run("default_level", func(t *testing.T) {
		if Level().Level() != globalLevel.Level() {
			t.Fatal("lever from context must be equal global level")
		}
	})

	t.Run("default_level_from_context", func(t *testing.T) {
		if LevelFromContext(context.Background()).Level() != globalLevel.Level() {
			t.Fatal("lever from context must be equal global level")
		}
	})

	t.Run("level_to_context", func(t *testing.T) {
		ctx := LevelToContext(context.Background(), LevelDebug)
		if LevelFromContext(ctx).Level() != LevelDebug {
			t.Fatal("lever from context must be equal debug level")
		}
	})

	t.Run("debug", func(t *testing.T) {
		buf := &bytes.Buffer{}
		ctx := context.Background()
		l := NewLogger(buf, false)

		SetLogger(l)
		SetLevel(LevelInfo)

		Debug(ctx, "msg")
		if buf.String() != "" {
			t.Fatal("debug log unexpected on info level")
		}

		SetLevel(LevelDebug)

		Debug(ctx, "msg")
		want := "\"msg\":\"msg\""
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("not found message in buffer: \ngot: %s;\nwant: %s", got, want)
		}
		buf.Reset()

		Debugf(ctx, "msg_%d", 1)
		want = "\"msg\":\"msg_1\""
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("not found message in buffer: \ngot: %s;\nwant: %s", got, want)
		}
	})

	t.Run("info", func(t *testing.T) {
		buf := &bytes.Buffer{}
		ctx := context.Background()
		l := NewLogger(buf, false)

		SetLogger(l)
		SetLevel(LevelWarn)

		Info(ctx, "msg")
		if buf.String() != "" {
			t.Fatal("info log unexpected on warn level")
		}

		SetLevel(LevelInfo)

		Info(ctx, "msg")
		want := "\"msg\":\"msg\""
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("not found message in buffer: \ngot: %s;\nwant: %s", got, want)
		}
		buf.Reset()

		Infof(ctx, "msg_%d", 1)
		want = "\"msg\":\"msg_1\""
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("not found message in buffer: \ngot: %s;\nwant: %s", got, want)
		}
	})

	t.Run("warn", func(t *testing.T) {
		buf := &bytes.Buffer{}
		ctx := context.Background()
		l := NewLogger(buf, false)

		SetLogger(l)
		SetLevel(LevelError)

		Warn(ctx, "msg")
		if buf.String() != "" {
			t.Fatal("warn log unexpected on error level")
		}

		SetLevel(LevelWarn)

		Warn(ctx, "msg")
		want := "\"msg\":\"msg\""
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("not found message in buffer: \ngot: %s;\nwant: %s", got, want)
		}
		buf.Reset()

		Warnf(ctx, "msg_%d", 1)
		want = "\"msg\":\"msg_1\""
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("not found message in buffer: \ngot: %s;\nwant: %s", got, want)
		}
	})

	t.Run("error", func(t *testing.T) {
		buf := &bytes.Buffer{}
		ctx := context.Background()
		l := NewLogger(buf, false)

		SetLogger(l)
		SetLevel(LevelError)

		Error(ctx, "msg")
		want := "\"msg\":\"msg\""
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("not found message in buffer: \ngot: %s;\nwant: %s", got, want)
		}
		buf.Reset()

		Errorf(ctx, "msg_%d", 1)
		want = "\"msg\":\"msg_1\""
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("not found message in buffer: \ngot: %s;\nwant: %s", got, want)
		}
	})

	t.Run("fatal", func(t *testing.T) {
		buf := &bytes.Buffer{}
		ctx := context.Background()
		l := NewLogger(buf, false)
		exitF = func(int) {}

		SetLogger(l)
		SetLevel(LevelError)

		Fatal(ctx, "msg")
		want := "\"msg\":\"msg\""
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("not found message in buffer: \ngot: %s;\nwant: %s", got, want)
		}
		buf.Reset()

		Fatalf(ctx, "msg_%d", 1)
		want = "\"msg\":\"msg_1\""
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("not found message in buffer: \ngot: %s;\nwant: %s", got, want)
		}
	})

	t.Run("with", func(t *testing.T) {
		ctx := context.Background()
		buf := &bytes.Buffer{}
		ctx = ToContext(ctx, NewLogger(buf, false))
		ctx = With(ctx, "arg", "val")
		want := "\"arg\":\"val\""

		Error(ctx, "test")
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("args in message: \ngot: %s;\nwant: %s", got, want)
		}
		buf.Reset()

		func(c context.Context) {
			With(c, "f", "f")
		}(ctx)

		Error(ctx, "test")
		want = "\"f\":\"f\""
		if got := buf.String(); strings.Contains(got, want) {
			t.Fatalf("args in message unexpected: \ngot: %s;\nwant: %s", got, want)
		}
	})

	t.Run("source", func(t *testing.T) {
		ctx := context.Background()
		buf := &bytes.Buffer{}
		ctx = ToContext(ctx, NewLogger(buf, true))
		want := "logger_test.go"

		Error(ctx, "test")
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("no args in message: \ngot: %s;\nwant: %s", got, want)
		}
		buf.Reset()

		FromContext(ctx).Error("test")
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("no args in message: \ngot: %s;\nwant: %s", got, want)
		}
		buf.Reset()
	})

	t.Run("with_group", func(t *testing.T) {
		ctx := context.Background()
		buf := &bytes.Buffer{}
		ctx = ToContext(ctx, NewLogger(buf, false))
		want := "\"group\":{\"key\":\"val\"}"

		WithGroup(ctx, "group").LogAttrs(ctx, LevelError, "msg", slog.Any("key", "val"))
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("no group in message: \ngot: %s;\nwant: %s", got, want)
		}
		buf.Reset()
	})
}

func Test_replaceSource(t *testing.T) {
	t.Run("has_no_source_key", func(t *testing.T) {
		attr := slog.Any("key", "val")
		if got := replaceSource(nil, attr); !reflect.DeepEqual(attr, got) {
			t.Fatalf("replace source: \ngot: %s;\nwant: %s", got, attr)
		}
	})
	t.Run("has_no_source_key", func(t *testing.T) {
		attr := slog.Any(slog.SourceKey, &slog.Source{Function: "Function", File: "File", Line: 1})
		want := slog.Any(slog.SourceKey, &slog.Source{Function: "", File: "File", Line: 1})
		if got := replaceSource(nil, attr); !reflect.DeepEqual(want, got) {
			t.Fatalf("replace source: \ngot: %s;\nwant: %s", got, want)
		}
	})
}
