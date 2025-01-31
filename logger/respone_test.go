package logger

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func removeAnsiCodes(str string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(str, "")
}

func TestLoggingResponse(t *testing.T) {
	var buf bytes.Buffer
	consoleLogger := log.New(&buf, "", log.LstdFlags)

	respLogger := NewResponLogger(consoleLogger)

	handler := respLogger.LoggingResponse(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	logOutput := removeAnsiCodes(buf.String())

	if !strings.Contains(logOutput, "200 OK") {
		t.Errorf("expected log to contain status code 200 OK, got %s", logOutput)
	}

	if !strings.Contains(logOutput, "trace id:") {
		t.Errorf("expected log to contain trace id, got %s", logOutput)
	}

	if !strings.Contains(logOutput, "url: /test") {
		t.Errorf("expected log to contain request URL '/test', got %s", logOutput)
	}

	if !strings.Contains(logOutput, "method: GET") {
		t.Errorf("expected log to contain method GET, got %s", logOutput)
	}

	xRequestID := rr.Header().Get("X-Request-ID")
	if xRequestID == "" {
		t.Errorf("expected X-Request-ID to be set in response header, but got empty")
	}

	if req.Header.Get("X-Request-ID") == "" {
		t.Errorf("expected X-Request-ID to be set in request header, but got empty")
	}
}
