package middleware

import (
	"net/http"

	"github.com/keenywheels/go-spy/pkg/httputils"
	"github.com/keenywheels/go-spy/pkg/logger"
)

// WithRecover middleware which provides recovery for panic
func WithRecover(l logger.Logger, next http.Handler) http.Handler {
	prefix := "RECOVER"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				l.Errorf("[%s] got err: %v", prefix, err)
				httputils.InternalErrorJSON(w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
