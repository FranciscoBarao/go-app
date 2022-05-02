package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
)

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	if r.Header.Get("Content-Type") != "" { // Only allow requests with application/json as header
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return NewError(http.StatusUnsupportedMediaType, msg)
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // Dont allow bodies that are over 1MB

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Dont allow any extra unexpected fields in the JSON

	err := decoder.Decode(&dst)
	if err == nil {
		err = decoder.Decode(&struct{}{})
		if err != io.EOF { // Don't allow several JSON objects
			msg := "Request body must only contain a single JSON object"
			return NewError(http.StatusBadRequest, msg)
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

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg = fmt.Sprintf("Request body contains unknown field %s", fieldName)

	case errors.Is(err, io.EOF):
		msg = "Request body must not be empty"

	case err.Error() == "http: request body too large":
		msg = "Request body must not be larger than 1MB"
		return NewError(http.StatusRequestEntityTooLarge, msg)

	default:
		return err
	}

	return NewError(http.StatusBadRequest, msg)
}
