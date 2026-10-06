package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/FranciscoBarao/catalog/internal/logging"
)

// DecodeJSONBody decodes a JSON request body into the given destination struct.
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	log := logging.FromCtx(r.Context())

	// Only allow requests with application/json as header
	mediaType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mediaType != "application/json" {
		return fmt.Errorf("content-type header is not application/json: %w", ErrBadRequest)
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // Dont allow bodies that are over 1MB

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Dont allow any extra unexpected fields in the JSON

	err := decoder.Decode(dst)
	if err == nil {
		if err = decoder.Decode(&struct{}{}); err != io.EOF { // Don't allow several JSON objects
			return fmt.Errorf("request body must only contain a single JSON object: %w", ErrBadRequest)
		}
		return nil // Success exit
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
		return fmt.Errorf("request body must not be larger than 1MB: %w", ErrRequestTooLarge)

	default:
		log.Error().Err(err).Msg("decode json body")
		return errors.New("internal error")
	}
	return fmt.Errorf("%s: %w", msg, ErrBadRequest)
}
