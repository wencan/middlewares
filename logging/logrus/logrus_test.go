package logrus_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sirupsen/logrus/hooks/test"

	"github.com/sirupsen/logrus"
	middleware_logrus "github.com/wencan/middlewares/logging/logrus"
)

func TestLoggingLogrus(t *testing.T) {
	logrusLogger := logrus.New()
	hook := test.NewLocal(logrusLogger)

	logger := middleware_logrus.NewLogger(logrusLogger)

	testMethod := http.MethodGet
	testURI := "/"
	testBodyByesSend := 123
	testStatus := http.StatusNotFound
	testTimestamp := time.Now()
	req := httptest.NewRequest(testMethod, testURI, nil)

	err := logger.Write(req, testStatus, testBodyByesSend, testTimestamp)
	assert.NoError(t, err)

	entries := hook.AllEntries()

	if assert.Equal(t, 1, len(entries)) {
		entry := entries[0]
		fields := entry.Data
		assert.Equal(t, testTimestamp, fields["time_local"])
		assert.Equal(t, testMethod, fields["request_method"])
		assert.Equal(t, testURI, fields["request_uri"])
		assert.Equal(t, testStatus, fields["status"])
		assert.Equal(t, testBodyByesSend, fields["body_bytes_sent"])
	}
}
