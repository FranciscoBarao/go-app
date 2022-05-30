package utils

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/unrolled/render"
)

func HTTPHandler(w http.ResponseWriter, returnValue interface{}, status int, err error) {
	if err != nil {
		var mr *MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	render.New().JSON(w, status, returnValue)
}

func GetFieldFromURL(r *http.Request, field string) string {
	return chi.URLParam(r, field)
}
