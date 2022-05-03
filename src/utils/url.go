package utils

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetFieldFromURL(r *http.Request, field string) string {
	return chi.URLParam(r, field)
}
