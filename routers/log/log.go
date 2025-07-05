package log

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/light-speak/lighthouse/logs"
)

type LogMiddleware struct{}

func (l *LogMiddleware) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &CustomLogEntry{Request: r}
}

type CustomLogEntry struct {
	Request *http.Request
}

func (c *CustomLogEntry) Write(status int, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	logger := logs.Debug()
	if status >= 400 {
		logger = logs.Error()
	}
	logger.
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Int("status", status).
		Int("bytes", bytes).
		Dur("elapsed", elapsed).
		Str("ip", c.Request.RemoteAddr).
		Interface("extra", extra).
		Msg("request")
}

func (c *CustomLogEntry) Panic(v interface{}, stack []byte) {
	logs.Error().Interface("panic", v).Bytes("stack", stack).Msg("panic")
}
