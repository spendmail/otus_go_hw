package internalhttp

import (
	"net/http"
	"time"

	internallogger "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(next http.HandlerFunc, logger internallogger.Interface) http.HandlerFunc {
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
