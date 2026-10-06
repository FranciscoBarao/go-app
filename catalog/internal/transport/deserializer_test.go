package transport

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type jsonBody struct {
	Name string `json:"name"`
	N    int    `json:"n"`
}

func TestDecodeJSONBody(t *testing.T) {
	tests := []struct {
		name         string
		contentType  string
		body         io.Reader
		want         jsonBody
		wantErr      error
		wantContains string
	}{
		{
			name:        "decodes one json object",
			contentType: "application/json",
			body:        strings.NewReader(`{"name":"Catan","n":2}`),
			want:        jsonBody{Name: "Catan", N: 2},
		},
		{
			name:        "accepts a charset parameter",
			contentType: "application/json; charset=utf-8",
			body:        strings.NewReader(`{"name":"Catan","n":2}`),
			want:        jsonBody{Name: "Catan", N: 2},
		},
		{
			name:         "rejects a non-json content type",
			contentType:  "text/plain",
			body:         strings.NewReader(`{"name":"Catan","n":2}`),
			wantErr:      utils.ErrBadRequest,
			wantContains: "content-type header is not application/json",
		},
		{
			name:         "rejects a missing content type",
			body:         strings.NewReader(`{"name":"Catan","n":2}`),
			wantErr:      utils.ErrBadRequest,
			wantContains: "content-type header is not application/json",
		},
		{
			name:         "rejects an empty body",
			contentType:  "application/json",
			body:         strings.NewReader(""),
			wantErr:      utils.ErrBadRequest,
			wantContains: "Request body must not be empty",
		},
		{
			name:         "rejects badly formed json",
			contentType:  "application/json",
			body:         strings.NewReader(`{`),
			wantErr:      utils.ErrBadRequest,
			wantContains: "badly-formed JSON",
		},
		{
			name:         "rejects a truncated json value",
			contentType:  "application/json",
			body:         strings.NewReader(`{"name":`),
			wantErr:      utils.ErrBadRequest,
			wantContains: "badly-formed JSON",
		},
		{
			name:         "rejects a value of the wrong type",
			contentType:  "application/json",
			body:         strings.NewReader(`{"name":"Catan","n":"two"}`),
			wantErr:      utils.ErrBadRequest,
			wantContains: `invalid value for the "n" field`,
		},
		{
			name:         "rejects an unknown field",
			contentType:  "application/json",
			body:         strings.NewReader(`{"name":"Catan","extra":1}`),
			wantErr:      utils.ErrBadRequest,
			wantContains: "unknown field",
		},
		{
			name:         "rejects more than one json object",
			contentType:  "application/json",
			body:         strings.NewReader(`{"name":"Catan","n":1}{"name":"Pandemic","n":2}`),
			wantErr:      utils.ErrBadRequest,
			wantContains: "single JSON object",
		},
		{
			name:         "rejects a body over 1MB",
			contentType:  "application/json",
			body:         strings.NewReader(`{"name":"` + strings.Repeat("a", int(maxJSONBodyBytes)) + `"}`),
			wantErr:      utils.ErrRequestTooLarge,
			wantContains: "must not be larger than 1MB",
		},
		{
			name:         "replaces an unexpected read error",
			contentType:  "application/json",
			body:         errReader{},
			wantContains: "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", tt.body)
			if tt.contentType != "" {
				request.Header.Set("Content-Type", tt.contentType)
			}

			var got jsonBody
			err := decodeJSONBody(httptest.NewRecorder(), request, &got)
			if tt.wantErr == nil && tt.wantContains == "" {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
				return
			}
			require.Error(t, err)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}
			assert.ErrorContains(t, err, tt.wantContains)
		})
	}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}
