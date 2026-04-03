package middleware

import (
	"errors"
	"net/http"
)

// ErrorHandler writes an appropriate HTTP error response based on the error type.
func ErrorHandler(w http.ResponseWriter, err error) {
	if err != nil {
		var mr *MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
