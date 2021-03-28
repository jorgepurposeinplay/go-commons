package middleware

import (
	"fmt"
	cmiddleware "github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)


func NewLoggerMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return cmiddleware.RequestLogger(&structuredLogger{logger})
}

type structuredLogger struct {
	Logger *zap.Logger
}

func (l *structuredLogger) NewLogEntry(r *http.Request) cmiddleware.LogEntry {
	entry := &structuredLoggerEntry{Logger: l.Logger}

	fields := []zapcore.Field{zap.String("ts", time.Now().UTC().Format(time.RFC1123))}

	if reqID := cmiddleware.GetReqID(r.Context()); reqID != "" {
		fields = append(fields, zap.String("req.id", reqID))
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	fields = append(fields, []zapcore.Field{
		zap.String("http_scheme", scheme),
		zap.String("http_proto", r.Proto),
		zap.String("http_method", r.Method),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()),
		zap.String("uri", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)),
	}...)

	entry.Logger = l.Logger.With(fields...)

	entry.Logger.Info("request started")

	return entry
}

type structuredLoggerEntry struct {
	Logger *zap.Logger
}

func (l *structuredLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.Logger = l.Logger.With(
		zap.Int("status", status),
		zap.Int("bytes_length", bytes),
		zap.Float64("duration_ms", float64(elapsed.Nanoseconds())/1000000.0),
	)

	l.Logger.Info("request complete")
}

func (l *structuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.With(
		zap.String("stack", string(stack)),
		zap.String("panic", fmt.Sprintf("%+v", v)),
	)
}

func GetLogEntry(r *http.Request) *zap.Logger {
	entry, _ := cmiddleware.GetLogEntry(r).(*structuredLoggerEntry)

	if entry == nil {
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		return logger
	}

	return entry.Logger
}
