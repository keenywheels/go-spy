package security

import (
	"context"
	"errors"

	gen "github.com/keenywheels/go-spy/internal/api/v1"
)

// ErrWrongToken is returned when token is wrong
var ErrWrongToken = errors.New("wrong api token")

// ctxKeyClient type for context key
type ctxKeyClient int

// clientKey value to put and get client from context
const clientKey ctxKeyClient = 0

// HandleS2STokenAuth handles S2STokenAuth security scheme
func (c *Controller) HandleS2STokenAuth(
	ctx context.Context,
	operationName gen.OperationName,
	t gen.S2STokenAuth,
) (context.Context, error) {
	gotToken := t.GetAPIKey()
	if gotToken == "" {
		return ctx, ErrWrongToken
	}

	for client, token := range c.clients {
		if token == gotToken {
			return context.WithValue(ctx, clientKey, client), nil
		}
	}

	return ctx, ErrWrongToken
}

// GetClientFromContext gets client from context
func GetClientFromContext(ctx context.Context) string {
	if client, ok := ctx.Value(clientKey).(string); ok {
		return client
	}

	return ""
}
