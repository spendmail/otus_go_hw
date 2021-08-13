package internalhttp

import (
	"net/http"
	"time"
)

type statusCodeWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusCodeWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Wraps request handler, adding request access log.
func loggingMiddleware(next http.HandlerFunc, logger Logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		wrapper := &statusCodeWrapper{writer, http.StatusOK}

		next(wrapper, request)

		logger.Info(
			"access",
			"addr", request.RemoteAddr,
			"method", request.Method,
			"path", request.URL.Path,
			"proto", request.Proto,
			"code", wrapper.statusCode,
			"latency", time.Since(start),
			"agent", request.UserAgent(),
		)
	}
}
