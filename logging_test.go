package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/wencan/middlewares"
	"github.com/wencan/middlewares/mock_logging"
)

func TestLoggingMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mocked_logger := mock_logging.NewMockLoggingLogger(ctrl)
	var received_req *http.Request
	var received_status, received_bodyBytesSend int
	mocked_logger.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(req *http.Request, status, bodyBytesSent int, timestamp time.Time) {
		received_req = req
		received_status = status
		received_bodyBytesSend = bodyBytesSent
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	recorder := httptest.NewRecorder()
	middleware := middlewares.LoggingMiddleware(mocked_logger)
	bodyString := []byte("hello, world")
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(bodyString)
	}))
	handler.ServeHTTP(recorder, req)

	assert.EqualValues(t, req, received_req)
	assert.Equal(t, http.StatusOK, received_status)
	assert.Equal(t, len(bodyString), received_bodyBytesSend)
}
