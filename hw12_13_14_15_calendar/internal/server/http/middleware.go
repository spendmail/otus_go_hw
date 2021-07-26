package internalhttp

import (
	"net/http"
	"time"
)

// Wraps request handler, adding request access log.
func loggingMiddleware(next http.HandlerFunc, logger Logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		next(writer, request)

		logger.Info(
			"access",
			"addr", request.RemoteAddr,
			"method", request.Method,
			"path", request.URL.Path,
			"proto", request.Proto,
			"code", 200,
			"latency", time.Since(start),
			"agent", request.UserAgent(),
		)
	}
}
