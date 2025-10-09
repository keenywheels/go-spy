package ctxutils

import (
	"context"

	"github.com/keenywheels/go-spy/pkg/logger"
)

// ctxKeyLogger using to set up logger in context
type ctxKeyLogger int

// loggerIDKey key to get/set data in context
const loggerIDKey ctxKeyLogger = 0

// SetLogger returns new context with logger
func SetLogger(ctx context.Context, logger logger.Logger) context.Context {
	return context.WithValue(ctx, loggerIDKey, logger)
}

// GetLogger return logger from context
func GetLogger(ctx context.Context) logger.Logger {
	if logger, ok := ctx.Value(loggerIDKey).(logger.Logger); ok {
		return logger
	}

	return nil
}
