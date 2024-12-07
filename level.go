package logger

import (
	"log/slog"
	"sync/atomic"
)

type LogLevel = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

var globalLevel = &atomicLevel{}

type atomicLevel struct {
	l atomic.Int32
}

func (al *atomicLevel) Level() LogLevel {
	return LogLevel(al.l.Load())
}

func (al *atomicLevel) SetLevel(l LogLevel) {
	al.l.Store(int32(l))
}
