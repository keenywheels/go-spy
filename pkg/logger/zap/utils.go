package zap

import (
	"go.uber.org/zap/zapcore"
)

// map to map cfg loglvl to zap loglvl
var loggerLevels = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// getLogLvl mapper cfg loglvl to zap loglvl
func (l *Logger) getLogLvl() zapcore.Level {
	lvl, ok := loggerLevels[l.loglvl]
	if !ok {
		return zapcore.DebugLevel
	}

	return lvl
}
