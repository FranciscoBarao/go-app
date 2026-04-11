package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/FranciscoBarao/catalog/boardgame"
	"github.com/FranciscoBarao/catalog/category"
	"github.com/FranciscoBarao/catalog/mechanism"
	"github.com/FranciscoBarao/catalog/middleware"
	tag "github.com/FranciscoBarao/catalog/tag"
)

type BoardGameSuite struct {
	suite.Suite
	base *Base
}

func (suite *BoardGameSuite) SetupSuite() {
	suite.base = NewBase(suite.T())
}

func (suite *BoardGameSuite) TestPostBoardgameSuccess() {
	bg := &boardgame.Boardgame{Name: "test", Publisher: "test", PlayerNumber: 1}
	suite.base.dbMock.EXPECT().
		CreateBoardgame(gomock.Any(), bg).
		Return(nil)

	bgJSON, err := json.Marshal(bg)
	suite.Require().NoError(err)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(bgJSON).
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *BoardGameSuite) TestPostExpansion() {
	// Expansion read of parent boardgame Mock
	parentIDStr := "1"
	parentID := uint(1)
	parentBg := boardgame.Boardgame{}
	parentBg.ID = parentID
	suite.base.dbMock.EXPECT().
		GetBoardgameByID(gomock.Any(), parentID).
		Return(parentBg, nil)

	// Boardgame expansion creation Mock
	expansion := &boardgame.Boardgame{Name: "expansion", Publisher: "expansion", PlayerNumber: 1}
	expansion.BoardgameID = &parentID
	suite.base.dbMock.EXPECT().
		CreateBoardgame(gomock.Any(), expansion).
		Return(nil)
	expansionJSON, err := json.Marshal(expansion)
	suite.Require().NoError(err)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame/"+parentIDStr+"/expansion").
		JSON(expansionJSON).
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *BoardGameSuite) TestGetBoardgame() {
	bgID := "1"
	expected := boardgame.Boardgame{}
	suite.base.dbMock.EXPECT().
		GetBoardgameByID(gomock.Any(), uint(1)).
		Return(expected, nil)

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
	suite.base.dbMock.EXPECT().
		DeleteBoardgame(gomock.Any(), uint(1)).
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
	suite.base.dbMock.EXPECT().
		GetTag(gomock.Any(), tagName).
		Return(tag.Tag{}, middleware.NewError(http.StatusNotFound, "Record not found"))

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
	suite.base.dbMock.EXPECT().
		GetCategory(gomock.Any(), categoryName).
		Return(category.Category{}, middleware.NewError(http.StatusNotFound, "Record not found"))

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
	suite.base.dbMock.EXPECT().
		GetMechanism(gomock.Any(), mechName).
		Return(mechanism.Mechanism{}, middleware.NewError(http.StatusNotFound, "Record not found"))

	apitest.New(). // Invalid Struct -> Mechanism does not previously exist
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[{"name":"test"}]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusNotFound).
			End()
	apitest.New(). // Invalid Struct -> Mechanisms have too many fields
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
	suite.base.dbMock.EXPECT().
		GetBoardgameByID(gomock.Any(), uint(1)).
		Return(boardgame.Boardgame{}, middleware.NewError(http.StatusNotFound, "Record not found"))

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
	suite.base.dbMock.EXPECT().
		DeleteBoardgame(gomock.Any(), uint(1)).
		Return(middleware.NewError(http.StatusNotFound, "Record not found"))

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
