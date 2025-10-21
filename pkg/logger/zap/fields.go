package zap

import (
	"github.com/keenywheels/go-spy/pkg/logger"
	"go.uber.org/zap"
)

// With method to add fields to logger
func (l *Logger) With(fields ...logger.Field) logger.Logger {
	if len(fields) == 0 {
		return l
	}

	// create a copy of the logger
	loggerWithFields := *l

	for _, field := range fields {
		loggerWithFields.sl = loggerWithFields.sl.With(
			zap.Any(field.Key, field.Value),
		)
	}

	return &loggerWithFields
}
