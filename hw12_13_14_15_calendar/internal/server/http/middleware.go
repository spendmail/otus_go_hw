package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.HandlerFunc, logger Logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		next(writer, request)

		logger.Info(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%d\t%s\t%s", request.RemoteAddr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			request.Method,
			request.URL.Path,
			request.Proto,
			200,
			time.Since(start),
			request.UserAgent(),
		))
	}
}
