package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"

	"github.com/FranciscoBarao/marketplace/internal/logging"
	"github.com/FranciscoBarao/marketplace/internal/middleware"
)

// DecodeJSONBody decodes a JSON request body into the given destination struct.
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	log := logging.FromCtx(r.Context())
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			log.Error().Str("content-type", value).Msg("content-type header must be application/json")
			return middleware.NewError(http.StatusBadRequest, "Content-Type header is not application/json")
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err == nil {
		if err = decoder.Decode(&struct{}{}); err != io.EOF {
			log.Error().Msg("request body must only contain a single json object")
			return middleware.NewError(http.StatusBadRequest, "Request body must only contain a single JSON object")
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
		return middleware.NewError(http.StatusRequestEntityTooLarge, "Request body must not be larger than 1MB")
	default:
		return err
	}
	return middleware.NewError(http.StatusBadRequest, msg)
}
