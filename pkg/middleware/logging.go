package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/keenywheels/go-spy/pkg/ctxutils"
	"github.com/keenywheels/go-spy/pkg/logger"
)

// responseData additional request data to be saved
type responseData struct {
	status int
	size   int
}

// loggingResponseWriter wrapper over http.ResponseWriter
type loggingResponseWriter struct {
	http.ResponseWriter
	data *responseData
}

// Write wrapper over http.ResponseWriter which saving size
func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.data.size += size

	return size, err
}

// WriteHeader wrapper over http.ResponseWriter which status code
func (w *loggingResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.data.status = status
}

// WithLogging logging incoming requests, also adds logger and reqid in request's context
func WithLogging(l logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			url    = r.URL.String()
			ip     = r.RemoteAddr
			method = r.Method
			start  = time.Now()
			reqid  = uuid.New().String()
		)

		l = l.With(logger.Field{Key: "reqid", Value: reqid})
		l.Infof("got request: url=%s, method=%s, ip=%s", url, method, ip)

		// create custom ResponseWriter
		lw := &loggingResponseWriter{
			ResponseWriter: w,
			data: &responseData{
				status: http.StatusOK,
				size:   0,
			},
		}

		// add reqid and logger in context
		ctx := ctxutils.SetRequestID(ctxutils.SetLogger(r.Context(), l), reqid)
		next.ServeHTTP(lw, r.WithContext(ctx))

		duration := time.Since(start)
		l.Infof("request to [%s] %s done: status=%s size=%d, time=%s",
			method,
			url,
			lw.data.status,
			lw.data.size,
			duration,
		)
	})
}
