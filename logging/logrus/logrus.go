package logrus

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

/*
 * Integrated logrus
 *
 * wencan
 * 2019-07-05
 */

// Logger Logrus logger. Implement of LoggingLogger
type Logger struct {
	logger *logrus.Logger
}

// NewLogger create logrus logger.
func NewLogger(logger *logrus.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func buildLogrusFields(req *http.Request, status, bodyBytesSent int, timestamp time.Time) logrus.Fields {
	fields := logrus.Fields{
		"remote_addr":     req.RemoteAddr,
		"time_local":      timestamp,
		"request_method":  req.Method,
		"request_uri":     req.RequestURI,
		"server_protocol": req.Proto,
		"status":          status,
		"body_bytes_sent": bodyBytesSent,
		"elapsed_time":    time.Now().Sub(timestamp).String(),
	}
	user, _, ok := req.BasicAuth()
	if ok {
		fields["remote_user"] = user
	}
	if req.Referer() != "" {
		fields["http_referer"] = req.Referer()
	}
	if req.UserAgent() != "" {
		fields["http_user_agent"] = req.UserAgent()
	}
	return fields
}

// Write logs a message.
func (logging *Logger) Write(req *http.Request, status, bodyBytesSent int, timestamp time.Time) error {
	fields := buildLogrusFields(req, status, bodyBytesSent, timestamp)
	logger := logging.logger.WithFields(fields)

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
