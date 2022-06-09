package utils

import (
	"catalog/model"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
)

func TestGetFiltersSuccess(t *testing.T) {
	apitest.New(). // name.a -> Names that contain letter a
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, value, err := GetFilters(model.Boardgame{}, "name.a")
			if err != nil || body != "name LIKE ?" || value != "%a%" {
				w.WriteHeader(http.StatusBadRequest)
			}

			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Expect(t).
		Status(http.StatusOK).
		End()

	apitest.New(). // price.lt.11 -> Prices lower than 11
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, value, err := GetFilters(model.Boardgame{}, "price.lt.11")
			if err != nil || body != "price < ?" || value != "11" {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Expect(t).
		Status(http.StatusOK).
		End()

	apitest.New(). // No filter
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, value, err := GetFilters(model.Boardgame{}, "")
			if err != nil || body != "" || value != "" {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestFiltersFailure(t *testing.T) {

	apitest.New(). // Different number of allowed Fields -> a.a.a.a || a
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, err := GetFilters(model.Boardgame{}, "name.a.a.a")
			_, _, err2 := GetFilters(model.Boardgame{}, "name")
			if err != nil && err2 != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}

			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Expect(t).
		Status(http.StatusUnprocessableEntity).
		End()

	apitest.New(). // Filters cant be empty -> .a || a..a || a.
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, err := GetFilters(model.Boardgame{}, ".name")
			_, _, err2 := GetFilters(model.Boardgame{}, "name.")
			_, _, err3 := GetFilters(model.Boardgame{}, "name..a")
			if err != nil && err2 != nil && err3 != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}

			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Expect(t).
		Status(http.StatusUnprocessableEntity).
		End()

	apitest.New(). // Operators must be "lt" || "le"|| "gt" || "ge" || "eq"
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, err := GetFilters(model.Boardgame{}, "price.asd.10")
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}

			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Expect(t).
		Status(http.StatusUnprocessableEntity).
		End()

	apitest.New(). // Fields must exist on Struct and be of the correct type (In this case Boardgame)
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, err := GetFilters(model.Boardgame{}, "test.a")       // Unknown field
			_, _, err2 := GetFilters(model.Boardgame{}, "price.lt.a")  // Incorrect type
			_, _, err3 := GetFilters(model.Boardgame{}, "name.eq.a")   // Incorrect type
			_, _, err4 := GetFilters(model.Boardgame{}, "price.10")    // Incorrect type
			_, _, err5 := GetFilters(model.Boardgame{}, "name.test_a") // Incorrect type
			if err != nil && err2 != nil && err3 != nil && err4 != nil && err5 != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Expect(t).
		Status(http.StatusUnprocessableEntity).
		End()
}

func TestGetSortsSuccess(t *testing.T) {
	apitest.New(). //
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sort, err := GetSort(model.Boardgame{}, "name.asc")
			if err != nil || sort != "name asc" {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Expect(t).
		Status(http.StatusOK).
		End()

	apitest.New(). //
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sort, err := GetSort(model.Boardgame{}, "price.desc")
			if err != nil || sort != "price desc" {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Expect(t).
		Status(http.StatusOK).
		End()

	apitest.New(). // No Sort
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sort, err := GetSort(model.Boardgame{}, "")
			if err != nil || sort != "" {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestSortsFailure(t *testing.T) {

	apitest.New(). // Different number of allowed Fields (2) -> a.a.a || a
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := GetSort(model.Boardgame{}, "name.asc.asc")
			_, err2 := GetSort(model.Boardgame{}, "name")
			if err != nil || err2 != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Expect(t).
		Status(http.StatusUnprocessableEntity).
		End()

	apitest.New(). // Filters cant be empty -> .a || a.
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := GetSort(model.Boardgame{}, ".asc")
			_, err2 := GetSort(model.Boardgame{}, "name.")
			if err != nil && err2 != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}

			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Expect(t).
		Status(http.StatusUnprocessableEntity).
		End()

	apitest.New(). // Order must be asc or desc
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := GetSort(model.Boardgame{}, "name.asd")
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}

			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Expect(t).
		Status(http.StatusUnprocessableEntity).
		End()

	apitest.New(). // Sorts must exist on Struct and be of a sortable type
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := GetSort(model.Boardgame{}, "test.a")    // Unknown field
			_, err2 := GetSort(model.Boardgame{}, "tags.asc") // Unsortable field
			if err != nil && err2 != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Expect(t).
		Status(http.StatusUnprocessableEntity).
		End()
}
