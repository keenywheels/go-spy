package ctxutils

import (
	"context"
)

// ctxKeyUserInfo using to set up user in context
type ctxKeyUserInfo int

// userInfoIDKey key to get/set data in context
const userInfoIDKey ctxKeyUserInfo = 0

// UserInfo contains user's info which stores in context
type UserInfo struct {
	Username string
	Email    string
}

// SetUserInfo returns new context with user info
func SetUserInfo(ctx context.Context, uinfo *UserInfo) context.Context {
	return context.WithValue(ctx, userInfoIDKey, uinfo)
}

// GetUserInfo returns user info from context
func GetUserInfo(ctx context.Context) *UserInfo {
	if uinfo, ok := ctx.Value(userInfoIDKey).(*UserInfo); ok {
		return uinfo
	}

	return nil
}
