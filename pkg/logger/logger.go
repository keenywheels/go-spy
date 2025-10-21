package logger

import "log"

// Field field used to add key-value pairs to logs
type Field struct {
	Key   string
	Value any
}

// Logger interface
type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Panicf(format string, args ...any)
	Fatalf(format string, args ...any)

	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Panic(args ...any)
	Fatal(args ...any)

	With(fields ...Field) Logger

	ToStdLog() *log.Logger
	Close() error
}
