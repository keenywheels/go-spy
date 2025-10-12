package ctxutils

import (
	"context"
)

// ctxKeyReqid using to set up request id in context
type ctxKeyReqid int

// reqidIDKey key to get/set data in context
const reqidIDKey ctxKeyReqid = 0

// SetRequestID returns new context with request id
func SetRequestID(ctx context.Context, reqid string) context.Context {
	return context.WithValue(ctx, reqidIDKey, reqid)
}

// GetRequestID return request id from context
func GetRequestID(ctx context.Context) string {
	if reqid, ok := ctx.Value(reqidIDKey).(string); ok {
		return reqid
	}

	return ""
}
