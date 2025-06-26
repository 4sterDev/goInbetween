package goInbetween

import (
	"io"
	"net/http"
	"time"

	"github.com/4sterDev/loggerRoom"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

var logger = &loggerRoom.Logger{
	Colorize:        true,
	LogLevel:        loggerRoom.INFO,
	ShowLogLevel:    true,
	ShowLogTimeDate: true,
	ShowLogTimeTime: true,
	ShowTag:         true,
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		request, err := io.ReadAll(r.Body)

		if err != nil {
			msg := loggerRoom.LogMessage{
				Tag:     "MWARE",
				Message: []any{err},
			}
			logger.Error(msg)
		}

		msg := loggerRoom.LogMessage{
			Tag:     "MWARE",
			Message: []any{wrapped.statusCode, r.Method, r.URL.Path, string(request), time.Since(start)},
		}

		logger.Info(msg)
	})
}
