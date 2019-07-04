package zap

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Logger Zap logger. Implement of LoggingLogger.
type Logger struct {
	logger *zap.Logger
}

// NewLogger create Zap logger.
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func buildZapFields(req *http.Request, status, bodyBytesSent int, timestamp time.Time) []zap.Field {
	fields := []zap.Field{}
	fields = append(fields, zap.String("remote_addr", req.RemoteAddr))
	user, _, ok := req.BasicAuth()
	if ok {
		fields = append(fields, zap.String("remote_user", user))
	}
	fields = append(fields, zap.Time("time_local", timestamp))
	fields = append(fields, zap.String("request_method", req.Method))
	fields = append(fields, zap.String("request_uri", req.RequestURI))
	fields = append(fields, zap.String("server_protocol", req.Proto))
	fields = append(fields, zap.Int("status", status))
	fields = append(fields, zap.Int("body_bytes_sent", bodyBytesSent))
	if req.Referer() != "" {
		fields = append(fields, zap.String("http_referer", req.Referer()))
	}
	if req.UserAgent() != "" {
		fields = append(fields, zap.String("http_user_agent", req.UserAgent()))
	}
	fields = append(fields, zap.String("elapsed_time", time.Now().Sub(timestamp).String()))
	return fields
}

// Write logs a message.
func (logging *Logger) Write(req *http.Request, status, bodyBytesSent int, timestamp time.Time) error {
	fields := buildZapFields(req, status, bodyBytesSent, timestamp)
	logger := logging.logger.With(fields...)

	statusText := http.StatusText(status)
	switch status / 100 {
	case 1, 2, 3:
		logger.Info(statusText)
	case 4:
		logger.Warn(statusText)
	case 5:
		logger.Error(statusText)
	default:
		logger.Warn("Unknown HTTP status")
	}
	return nil
}
