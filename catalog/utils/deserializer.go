package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/middleware/logging"
)

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	log := logging.FromCtx(context.Background())
	if r.Header.Get("Content-Type") != "" { // Only allow requests with application/json as header
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			log.Error().Str("content-type", value).Msg("content-type header of request must be application/json")
			return middleware.NewError(http.StatusBadRequest, "Content-Type header is not application/json")
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // Dont allow bodies that are over 1MB

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Dont allow any extra unexpected fields in the JSON

	err := decoder.Decode(&dst)
	if err == nil {
		if err = decoder.Decode(&struct{}{}); err != io.EOF { // Don't allow several JSON objects
			log.Error().Err(err).Msg("request body must only contain a single json object")
			return middleware.NewError(http.StatusBadRequest, "Request body must only contain a single JSON object")
		}
	}

	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var msg string

	switch {
	case errors.As(err, &syntaxError):
		log.Error().Int64("position", syntaxError.Offset).Msg("request body contains badly-formed json")
		msg = fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

	case errors.Is(err, io.ErrUnexpectedEOF):
		log.Error().Msg("request body contains badly-formed json")
		msg = "Request body contains badly-formed JSON"

	case errors.As(err, &unmarshalTypeError):
		log.Error().Str("field", unmarshalTypeError.Field).Int64("position", unmarshalTypeError.Offset).Msg("request body contains invalid value in field")
		msg = fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)

	case strings.HasPrefix(err.Error(), "json: unknown field"):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		log.Error().Str("field_name", fieldName).Msg("request body contains unknown field")
		msg = fmt.Sprintf("Request body contains unknown field %s", fieldName)

	case errors.Is(err, io.EOF):
		log.Error().Msg("request body is empty")
		msg = "Request body must not be empty"

	case err.Error() == "http: request body too large":
		log.Error().Msg("request body must not be larger than 1MB")
		return middleware.NewError(http.StatusRequestEntityTooLarge, "Request body must not be larger than 1MB")

	default:
		log.Error().Err(err)
		return err
	}
	return middleware.NewError(http.StatusBadRequest, msg)
}
