package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/utils"
)

// maxJSONBodyBytes is the maximum size of a JSON body in bytes (1MB).
const maxJSONBodyBytes int64 = 1 << 20

// decodeJSONBody decodes one JSON object from the request into dst.
// The body must be application/json and at most 1MB.
func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	mediaType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mediaType != "application/json" {
		return fmt.Errorf("content-type header is not application/json: %w", utils.ErrBadRequest)
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err == nil {
		if err = decoder.Decode(&struct{}{}); err != io.EOF {
			return fmt.Errorf("request body must only contain a single JSON object: %w", utils.ErrBadRequest)
		}
		return nil
	}

	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var msg string

	switch {
	case errors.As(err, &syntaxError):
		msg = fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg = "Request body contains badly-formed JSON"
	case errors.As(err, &unmarshalTypeError):
		msg = fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
	case strings.HasPrefix(err.Error(), "json: unknown field"):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg = fmt.Sprintf("Request body contains unknown field %s", fieldName)
	case errors.Is(err, io.EOF):
		msg = "Request body must not be empty"
	case err.Error() == "http: request body too large":
		return fmt.Errorf("request body must not be larger than 1MB: %w", utils.ErrRequestTooLarge)
	default:
		logging.FromCtx(r.Context()).Error().Err(err).Msg("decode json body")
		return errors.New("internal error")
	}
	return fmt.Errorf("%s: %w", msg, utils.ErrBadRequest)
}
