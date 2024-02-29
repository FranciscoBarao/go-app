package tests

import (
	"net/http"
	"testing"

	"catalog/model"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
)

type BoardGameSuite struct {
	suite.Suite

	base *Base
	id   string // Used for BG expansions
}

func (suite *BoardGameSuite) SetupSuite() {
	suite.base = NewBase(suite.T())
}

func (suite *BoardGameSuite) TestCreateBoardgameSuccess(t *testing.T) {
	bg := model.NewBoardgame("test", "test", 1, nil, nil, nil)
	suite.base.dbMock.EXPECT().
		Create(bg, "").
		Return(bg, "")

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"Name":"test","Publisher":"test","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func (suite *BoardGameSuite) TestGetAllBoardgameSuccess(t *testing.T) {
	// Verifies if the created boardgame exists on GetAll
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/boardgame").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		End()
}
func (suite *BoardGameSuite) TestCreateExpansionSuccess(t *testing.T) {
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame/"+suite.id+"/expansion").
		JSON(`{"Name":"expansion","Publisher":"expansion","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func (suite *BoardGameSuite) TestGetBoardgameSuccess(t *testing.T) {
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/boardgame/"+suite.id).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func (suite *BoardGameSuite) TestDeleteBoardgameSuccess(t *testing.T) {
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/boardgame/"+suite.id).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusNoContent).
		End()
}

func (suite *BoardGameSuite) TestCreateBoardgameJsonFailures(t *testing.T) {
	// Several Json Objects on the body
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`[{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]},{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}]`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"Name:"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"Name":100,"Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"TEST":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(``).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()
}

func (suite *BoardGameSuite) TestCreateBoardgameValidStructFailures(t *testing.T) {
	//  <<<< field - Name >>>>
	apitest.New(). // Invalid Struct -> NOT maxstringlength(100)
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	apitest.New(). // Invalid Struct -> NOT alphanum
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test?","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	//  <<<< field - Publisher >>>>
	apitest.New(). // Invalid Struct -> NOT maxstringlength(100)
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	apitest.New(). // Invalid Struct -> NOT alphanum
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test?","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	//  <<<< field - Price >>>>
	apitest.New(). // Invalid Struct -> NOT float
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":"string","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

	apitest.New(). // Invalid Struct -> NOT in range(0|1000)
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":1001,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	//  <<<< field - PlayerNumber >>>>
	apitest.New(). // Invalid Struct -> NOT int
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1.5,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

	apitest.New(). // Invalid Struct -> NOT in range(0|16)
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":17,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	//  <<<< field - Tags >>>>
	apitest.New(). // Invalid Struct -> Tag does not previously exist
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[{"name":"test"}],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(t).
			Status(http.StatusNotFound).
			End()

	apitest.New(). // Invalid Struct -> Tags have too many fields
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[{"name":"test", "test":"test"}],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

	//  <<<< field - Categories >>>>
	apitest.New(). // Invalid Struct -> Category does not previously exist
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[{"name":"test"}],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	apitest.New(). // Invalid Struct -> Categories have too many fields
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[{"name":"test", "test":"test"}],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

	//  <<<< field - Mechanisms >>>>
	apitest.New(). // Invalid Struct -> Mechanism does not previously exist
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[{"name":"test"}]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	apitest.New(). // Invalid Struct -> Mchanisms have too many fields
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[{"name":"test", "test":"test"}]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
}

func (suite *BoardGameSuite) TestGetBoardgameFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/boardgame/1000").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
func (suite *BoardGameSuite) TestDeleteBoardgameFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/boardgame/1000").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestBoardGameSuite(t *testing.T) {
	suite.Run(t, new(BoardGameSuite))
}
