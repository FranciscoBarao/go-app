package tests

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/model"
)

type BoardGameSuite struct {
	suite.Suite
	base *Base
}

func (suite *BoardGameSuite) SetupSuite() {
	suite.base = NewBase(suite.T())
}

func (suite *BoardGameSuite) TestPostBoardgameSuccess() {
	bg := &model.Boardgame{Name: "test", Publisher: "test", PlayerNumber: 1}
	suite.base.dbMock.EXPECT().
		Create(bg).
		Return(nil)

	bgJson, err := json.Marshal(bg)
	suite.Require().NoError(err)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(bgJson).
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *BoardGameSuite) TestPostExpansion() {
	// Expansion read of parent boardgame Mock
	parentIDStr := "0"
	u64, _ := strconv.ParseUint(parentIDStr, 10, 32)
	parentID := uint(u64)
	parentBg := new(model.Boardgame)
	suite.base.dbMock.EXPECT().
		Read(parentBg, "", "id = ?", parentIDStr).
		Return(nil)

	// Boardgame expansion creation Mock
	expansion := &model.Boardgame{Name: "expansion", Publisher: "expansion", PlayerNumber: 1}
	expansion.SetBoardgameID(&parentID)
	suite.base.dbMock.EXPECT().
		Create(expansion).
		Return(nil)
	expansionJson, err := json.Marshal(expansion)
	suite.Require().NoError(err)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame/"+parentIDStr+"/expansion").
		JSON(expansionJson).
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *BoardGameSuite) TestGetBoardgame() {
	bgID := "test"
	bg := new(model.Boardgame)
	suite.base.dbMock.EXPECT().
		Read(bg, "", "id = ?", bgID).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/boardgame/"+bgID).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *BoardGameSuite) TestDeleteBoardgameSuccess() {
	bgID := "1"
	bg := new(model.Boardgame)
	suite.base.dbMock.EXPECT().
		Read(bg, "", "id = ?", bgID).
		Return(nil)

	suite.base.dbMock.EXPECT().
		Delete(new(model.Boardgame)).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/boardgame/"+bgID).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()
}

func (suite *BoardGameSuite) TestPostBoardgameJsonFailures() {
	// Several Json Objects on the body
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`[{"Name":"test","Publisher":"test","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]},{"Name":"test","Publisher":"test","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}]`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"Name:"test","Publisher":"test","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"Name":100,"Publisher":"test","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"TEST":"test","Publisher":"test","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(``).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *BoardGameSuite) TestPostBoardgameStructFailures() {
	//  <<<< field - Name >>>>
	apitest.New(). // Invalid Struct -> NOT maxstringlength(100)
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","Publisher":"test","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusForbidden).
			End()

	apitest.New(). // Invalid Struct -> NOT alphanum
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test?","Publisher":"test","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusForbidden).
			End()

	//  <<<< field - Publisher >>>>
	apitest.New(). // Invalid Struct -> NOT maxstringlength(100)
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusForbidden).
			End()

	apitest.New(). // Invalid Struct -> NOT alphanum
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test?","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusForbidden).
			End()

	//  <<<< field - PlayerNumber >>>>
	apitest.New(). // Invalid Struct -> NOT int
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","PlayerNumber":1.5,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusBadRequest).
			End()

	apitest.New(). // Invalid Struct -> NOT in range(0|16)
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","PlayerNumber":17,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusForbidden).
			End()
}

func (suite *BoardGameSuite) TestPostBoardgameAssociationFailures() {
	//  <<<< field - Tags >>>>
	tagName := "test"
	tag := new(model.Tag)
	suite.base.dbMock.EXPECT().
		Read(tag, "", "name = ?", tagName).
		Return(middleware.NewError(http.StatusNotFound, "Record not found"))

	apitest.New(). // Invalid Struct -> Tag does not previously exist
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","PlayerNumber":1,"Tags":[{"name":"`+tagName+`"}],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusNotFound).
			End()

	apitest.New(). // Invalid Struct -> Tags have too many fields
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","PlayerNumber":1,"Tags":[{"name":"test", "test":"test"}],"Categories":[],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusBadRequest).
			End()

	//  <<<< field - Categories >>>>
	categoryName := "test"
	category := new(model.Category)
	suite.base.dbMock.EXPECT().
		Read(category, "", "name = ?", categoryName).
		Return(middleware.NewError(http.StatusNotFound, "Record not found"))

	apitest.New(). // Invalid Struct -> Category does not previously exist
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","PlayerNumber":1,"Tags":[],"Categories":[{"name":"test"}],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusNotFound).
			End()
	apitest.New(). // Invalid Struct -> Categories have too many fields
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","PlayerNumber":1,"Tags":[],"Categories":[{"name":"test", "test":"test"}],"Mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusBadRequest).
			End()

	//  <<<< field - Mechanisms >>>>
	mechName := "test"
	mech := new(model.Mechanism)
	suite.base.dbMock.EXPECT().
		Read(mech, "", "name = ?", mechName).
		Return(middleware.NewError(http.StatusNotFound, "Record not found"))

	apitest.New(). // Invalid Struct -> Mechanism does not previously exist
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[{"name":"test"}]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusNotFound).
			End()
	apitest.New(). // Invalid Struct -> Mchanisms have too many fields
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[{"name":"test", "test":"test"}]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusBadRequest).
			End()
}

func (suite *BoardGameSuite) TestGetBoardgameFailure() {
	bgID := "1"
	bg := new(model.Boardgame)
	suite.base.dbMock.EXPECT().
		Read(bg, "", "id = ?", bgID).
		Return(middleware.NewError(http.StatusNotFound, "Boardgame not found with name: "+bgID))

	// Record not found
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/boardgame/"+bgID).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *BoardGameSuite) TestDeleteBoardgameFailure() {
	bgID := "1"
	bg := new(model.Boardgame)
	suite.base.dbMock.EXPECT().
		Read(bg, "", "id = ?", bgID).
		Return(middleware.NewError(http.StatusNotFound, "Boardgame not found with id: "+bgID))

	// Record not found
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/boardgame/"+bgID).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func TestBoardGameSuite(t *testing.T) {
	suite.Run(t, new(BoardGameSuite))
}
