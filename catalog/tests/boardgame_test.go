package tests

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"catalog/model"

	"github.com/steinfletcher/apitest"
)

var id uint

func TestCreateBoardgameSuccess(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		Assert(func(res *http.Response, req *http.Request) error { // Gets ID from BG Creation for further test use
			var bg model.Boardgame
			json.NewDecoder(res.Body).Decode(&bg)
			id = *bg.GetId()
			return nil
		}).
		End()
}

func TestGetAllBoardgameSuccess(t *testing.T) {
	// Verifies if the created boardgame exists on GetAll
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/boardgame").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		End()
}
func TestCreateExpansionSuccess(t *testing.T) {
	idString := strconv.FormatUint(uint64(id), 10)

	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/boardgame/"+idString+"/expansion").
		JSON(`{"Name":"expansion","Publisher":"expansion","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestGetBoardgameSuccess(t *testing.T) {
	idString := strconv.FormatUint(uint64(id), 10)
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/boardgame/"+idString).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestDeleteBoardgameSuccess(t *testing.T) {
	idString := strconv.FormatUint(uint64(id), 10)
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/boardgame/"+idString).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusNoContent).
		End()
}

func TestCreateBoardgameJsonFailures(t *testing.T) {
	// Several Json Objects on the body
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`[{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]},{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}]`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"Name:"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"Name":100,"Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"TEST":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/boardgame").
		JSON(``).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()
}

func TestCreateBoardgameValidStructFailures(t *testing.T) {
	//  <<<< field - Name >>>>
	apitest.New(). // Invalid Struct -> NOT maxstringlength(100)
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+oauthHeader).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	apitest.New(). // Invalid Struct -> NOT alphanum
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test?","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+oauthHeader).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	//  <<<< field - Publisher >>>>
	apitest.New(). // Invalid Struct -> NOT maxstringlength(100)
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+oauthHeader).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	apitest.New(). // Invalid Struct -> NOT alphanum
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test?","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+oauthHeader).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	//  <<<< field - Price >>>>
	apitest.New(). // Invalid Struct -> NOT float
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":"string","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+oauthHeader).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

	apitest.New(). // Invalid Struct -> NOT in range(0|1000)
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":1001,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+oauthHeader).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	//  <<<< field - PlayerNumber >>>>
	apitest.New(). // Invalid Struct -> NOT int
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1.5,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+oauthHeader).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

	apitest.New(). // Invalid Struct -> NOT in range(0|16)
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":17,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+oauthHeader).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	//  <<<< field - Tags >>>>
	apitest.New(). // Invalid Struct -> Tag does not previously exist
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[{"name":"test"}],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+oauthHeader).
			Expect(t).
			Status(http.StatusNotFound).
			End()

	apitest.New(). // Invalid Struct -> Tags have too many fields
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[{"name":"test", "test":"test"}],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+oauthHeader).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

	//  <<<< field - Categories >>>>
	apitest.New(). // Invalid Struct -> Category does not previously exist
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[{"name":"test"}],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+oauthHeader).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	apitest.New(). // Invalid Struct -> Categories have too many fields
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[{"name":"test", "test":"test"}],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+oauthHeader).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

	//  <<<< field - Mechanisms >>>>
	apitest.New(). // Invalid Struct -> Mechanism does not previously exist
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[{"name":"test"}]}`).
			Header("Authorization", "Bearer "+oauthHeader).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	apitest.New(). // Invalid Struct -> Mchanisms have too many fields
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[{"name":"test", "test":"test"}]}`).
			Header("Authorization", "Bearer "+oauthHeader).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
}

func TestGetBoardgameFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/boardgame/1000").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
func TestDeleteBoardgameFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/boardgame/1000").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
