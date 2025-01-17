package middleware

import (
	"log/slog"
	"net/http"
)

// Recover is a middleware that recover from panics that occur in the handlers,
// logs the error, and returns a 500 status code.
func Recover(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rc := recover(); rc != nil {

				wrapped := &wrappedWriter{
					ResponseWriter: w,
					statusCode:     500,
				}

				next.ServeHTTP(wrapped, r)

				logger.InfoContext(
					r.Context(),
					"panic recovered",
					slog.Any("error", rc),
					slog.Int("status", wrapped.statusCode),
				)
			}

		})
	}
}
