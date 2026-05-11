package tests

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/utils"
)

type UtilSuite struct {
	suite.Suite

	base *Base
}

func (suite *UtilSuite) SetupSuite() {
	suite.base = NewBase(suite.T())
}

func (suite *UtilSuite) TestGetFilters() {
	apitest.New(). // name.a -> Names that contain letter a
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, value, err := utils.GetFilters(boardgame.Boardgame{}, "name.a")
			if err != nil || body != "name LIKE ?" || value != "%a%" {
				w.WriteHeader(http.StatusBadRequest)
			}

			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()

	apitest.New(). // playernumber.lt.5 -> Player Number lower than 5
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, value, err := utils.GetFilters(boardgame.Boardgame{}, "playernumber.lt.5")
			if err != nil || body != "playernumber < ?" || value != "5" {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()

	apitest.New(). // No filter
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, value, err := utils.GetFilters(boardgame.Boardgame{}, "")
			if err != nil || body != "" || value != "" {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *UtilSuite) TestFiltersFailure() {
	apitest.New(). // Different number of allowed Fields -> a.a.a.a || a
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, err := utils.GetFilters(boardgame.Boardgame{}, "name.a.a.a")
			_, _, err2 := utils.GetFilters(boardgame.Boardgame{}, "name")
			if err != nil && err2 != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}

			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusUnprocessableEntity).
		End()

	apitest.New(). // Filters cant be empty -> .a || a..a || a.
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, err := utils.GetFilters(boardgame.Boardgame{}, ".name")
			_, _, err2 := utils.GetFilters(boardgame.Boardgame{}, "name.")
			_, _, err3 := utils.GetFilters(boardgame.Boardgame{}, "name..a")
			if err != nil && err2 != nil && err3 != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}

			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusUnprocessableEntity).
		End()

	apitest.New(). // Operators must be "lt" || "le"|| "gt" || "ge" || "eq"
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, err := utils.GetFilters(boardgame.Boardgame{}, "price.asd.10")
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}

			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusUnprocessableEntity).
		End()

	apitest.New(). // Fields must exist on Struct and be of the correct type (In this case Boardgame)
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, err := utils.GetFilters(boardgame.Boardgame{}, "test.a")       // Unknown field
			_, _, err2 := utils.GetFilters(boardgame.Boardgame{}, "price.lt.a")  // Incorrect type
			_, _, err3 := utils.GetFilters(boardgame.Boardgame{}, "name.eq.a")   // Incorrect type
			_, _, err4 := utils.GetFilters(boardgame.Boardgame{}, "price.10")    // Incorrect type
			_, _, err5 := utils.GetFilters(boardgame.Boardgame{}, "name.test_a") // Incorrect type
			if err != nil && err2 != nil && err3 != nil && err4 != nil && err5 != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusUnprocessableEntity).
		End()
}

func (suite *UtilSuite) TestGetSorts() {
	apitest.New(). //
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			col, order, err := utils.GetSort(boardgame.Boardgame{}, "name.asc")
			if err != nil || col != "name" || order != "asc" {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()

	apitest.New(). //
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			col, order, err := utils.GetSort(boardgame.Boardgame{}, "playernumber.desc")
			if err != nil || col != "player_number" || order != "desc" {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()

	apitest.New(). // No Sort
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			col, order, err := utils.GetSort(boardgame.Boardgame{}, "")
			if err != nil || col != "" || order != "" {
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *UtilSuite) TestSortsFailure() {
	apitest.New(). // Different number of allowed Fields (2) -> a.a.a || a
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, err := utils.GetSort(boardgame.Boardgame{}, "name.asc.asc")
			_, _, err2 := utils.GetSort(boardgame.Boardgame{}, "name")
			if err != nil || err2 != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusUnprocessableEntity).
		End()

	apitest.New(). // Filters cant be empty -> .a || a.
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, err := utils.GetSort(boardgame.Boardgame{}, ".asc")
			_, _, err2 := utils.GetSort(boardgame.Boardgame{}, "name.")
			if err != nil && err2 != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}

			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusUnprocessableEntity).
		End()

	apitest.New(). // Order must be asc or desc
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, err := utils.GetSort(boardgame.Boardgame{}, "name.asd")
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}

			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusUnprocessableEntity).
		End()

	apitest.New(). // Sorts must exist on Struct and be of a sortable type
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, err := utils.GetSort(boardgame.Boardgame{}, "test.asc")  // Unknown field
			_, _, err2 := utils.GetSort(boardgame.Boardgame{}, "tags.asc") // Unsortable field
			if err != nil && err2 != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusUnprocessableEntity).
		End()
}

func TestUtilSuite(t *testing.T) {
	suite.Run(t, new(UtilSuite))
}
