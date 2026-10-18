package logging

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit_UnknownLevelUsesDebug(t *testing.T) {
	orig := os.Stderr
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	os.Stderr = writer
	t.Cleanup(func() {
		os.Stderr = orig
	})

	Init("loud")
	require.NoError(t, writer.Close())

	out, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, zerolog.DebugLevel, zerolog.GlobalLevel())
	assert.Contains(t, string(out), "unknown log level")
	assert.Contains(t, string(out), "loud")
}

func TestRequestMiddleware_LogsRequestIDAndRoute(t *testing.T) {
	Init("info")
	var buf bytes.Buffer
	logger := zerolog.New(&buf)
	zerolog.DefaultContextLogger = &logger

	handler := RequestMiddleware(WithRequestIDHeader("X-Request-ID"))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/boardgame/catan?q=secret", nil)
	req.Header.Set("X-Request-ID", "A3F1C9E0-4B2D-4716-AA09-C4E71D55B083")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, "a3f1c9e0-4b2d-4716-aa09-c4e71d55b083", rec.Header().Get("X-Request-ID"))
	assert.Equal(t, http.StatusNoContent, rec.Code)

	var event map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &event))
	assert.Equal(t, "a3f1c9e0-4b2d-4716-aa09-c4e71d55b083", event["request_id"])
	assert.Equal(t, "/boardgame/catan", event["path"])
	assert.Equal(t, http.MethodGet, event["method"])
	assert.NotContains(t, buf.String(), "q=secret")
}

func TestRequestMiddleware_ReplacesUnsafeRequestID(t *testing.T) {
	Init("error")

	handler := RequestMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("X-Request-ID", "bad id\n")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	got := rec.Header().Get("X-Request-ID")
	_, err := uuid.Parse(got)
	assert.NotEqual(t, "bad id\n", got)
	assert.NoError(t, err)
}
