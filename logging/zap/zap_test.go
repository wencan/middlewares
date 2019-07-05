package zap_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	middleware_zap "github.com/wencan/middlewares/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingZap(t *testing.T) {
	core, observer := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core)

	logger := middleware_zap.NewLogger(zapLogger)

	testMethod := http.MethodGet
	testURI := "/"
	testBodyByesSend := 123
	testStatus := http.StatusNotFound
	testTimestamp := time.Now()
	req := httptest.NewRequest(testMethod, testURI, nil)

	err := logger.Write(req, testStatus, testBodyByesSend, testTimestamp)
	assert.NoError(t, err)

	entries := observer.All()

	if assert.Equal(t, 1, len(entries)) {
		entry := entries[0]
		fields := entry.ContextMap()
		assert.Equal(t, testTimestamp.Unix(), fields["time_local"].(time.Time).Unix())
		assert.Equal(t, testMethod, fields["request_method"])
		assert.Equal(t, testURI, fields["request_uri"])
		assert.EqualValues(t, testStatus, fields["status"])
		assert.EqualValues(t, testBodyByesSend, fields["body_bytes_sent"])
	}
}
