package tests

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/listopt"
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
	apitest.New(). // empty op -> partial match
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			col, op, val, numeric, err := utils.GetFilterFields(boardgame.Boardgame{}, "name", "", "a")
			if err != nil || col != "name" || op != "like" || val != "a" || numeric {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()

	apitest.New(). // lt -> numeric comparison
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			col, op, val, numeric, err := utils.GetFilterFields(boardgame.Boardgame{}, "maxplayers", "lt", "5")
			if err != nil || col != "max_players" || op != listopt.OpLt || val != "5" || !numeric {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()

	apitest.New(). // eq -> exact equality
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			col, op, val, numeric, err := utils.GetFilterFields(boardgame.Boardgame{}, "name", "eq", "Catan")
			if err != nil || col != "name" || op != "eq" || val != "Catan" || numeric {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()

	apitest.New(). // decimal value (dot in value) -> numeric comparison
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			col, op, val, numeric, err := utils.GetFilterFields(boardgame.Boardgame{}, "maxplayers", "ge", "3.5")
			if err != nil || col != "max_players" || op != listopt.OpGe || val != "3.5" || !numeric {
				w.WriteHeader(http.StatusBadRequest)
				return
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
	apitest.New(). // Empty field or value
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, _, _, err := utils.GetFilterFields(boardgame.Boardgame{}, "", "eq", "name")
			_, _, _, _, err2 := utils.GetFilterFields(boardgame.Boardgame{}, "name", "eq", "")
			if err != nil && err2 != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusUnprocessableEntity).
		End()

	apitest.New(). // Invalid operator
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, _, _, err := utils.GetFilterFields(boardgame.Boardgame{}, "playernumber", "asd", "10")
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			w.WriteHeader(http.StatusOK)
		}).
		Get("").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusUnprocessableEntity).
		End()

	apitest.New(). // Type mismatches
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, _, _, err := utils.GetFilterFields(boardgame.Boardgame{}, "unknown", "eq", "asc")        // Unknown field
			_, _, _, _, err2 := utils.GetFilterFields(boardgame.Boardgame{}, "playernumber", "lt", "abc") // Non-numeric value
			_, _, _, _, err3 := utils.GetFilterFields(boardgame.Boardgame{}, "name", "lt", "5")           // Numeric op on string
			_, _, _, _, err4 := utils.GetFilterFields(boardgame.Boardgame{}, "playernumber", "", "hello") // Like on non-string
			if err != nil && err2 != nil && err3 != nil && err4 != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
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
			col, order, err := utils.GetSort(boardgame.Boardgame{}, "maxplayers.desc")
			if err != nil || col != "max_players" || order != "desc" {
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
			_, _, err := utils.GetSort(boardgame.Boardgame{}, "test.asc")        // Unknown field
			_, _, err2 := utils.GetSort(boardgame.Boardgame{}, "mechanisms.asc") // Unsortable field
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
