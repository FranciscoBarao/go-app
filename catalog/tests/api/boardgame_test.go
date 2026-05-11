package tests

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	tag "github.com/FranciscoBarao/catalog/internal/tag"
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

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"name":"test","publisher":"test","playerNumber":1}`).
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

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame/"+parentIDStr+"/expansion").
		JSON(`{"name":"expansion","publisher":"expansion","playerNumber":1}`).
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
		JSON(`[{"name":"test","publisher":"test","playerNumber":1,"tags":[],"categories":[],"mechanisms":[]},{"name":"test","publisher":"test","playerNumber":1,"tags":[],"categories":[],"mechanisms":[]}]`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"name:"test","publisher":"test","playerNumber":1,"tags":[],"categories":[],"mechanisms":[]}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"name":100,"publisher":"test","playerNumber":1,"tags":[],"categories":[],"mechanisms":[]}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"TEST":"test","publisher":"test","playerNumber":1,"tags":[],"categories":[],"mechanisms":[]}`).
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
			JSON(`{"name":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","publisher":"test","playerNumber":1,"tags":[],"categories":[],"mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusBadRequest).
			End()

	//  <<<< field - Publisher >>>>
	apitest.New(). // Invalid Struct -> NOT maxstringlength(100)
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"name":"test","publisher":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","playerNumber":1,"tags":[],"categories":[],"mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusBadRequest).
			End()

	//  <<<< field - PlayerNumber >>>>
	apitest.New(). // Invalid Struct -> NOT int
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"name":"test","publisher":"test","playerNumber":1.5,"tags":[],"categories":[],"mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusBadRequest).
			End()

	apitest.New(). // Invalid Struct -> NOT in range(0|16)
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"name":"test","publisher":"test","playerNumber":17,"tags":[],"categories":[],"mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusBadRequest).
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
			JSON(`{"name":"test","publisher":"test","playerNumber":1,"tags":[{"name":"`+tagName+`"}],"categories":[],"mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusNotFound).
			End()

	apitest.New(). // Invalid Struct -> Tags have too many fields
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"name":"test","publisher":"test","playerNumber":1,"tags":[{"name":"test", "test":"test"}],"categories":[],"mechanisms":[]}`).
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
			JSON(`{"name":"test","publisher":"test","playerNumber":1,"tags":[],"categories":[{"name":"test"}],"mechanisms":[]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusNotFound).
			End()
	apitest.New(). // Invalid Struct -> Categories have too many fields
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"name":"test","publisher":"test","playerNumber":1,"tags":[],"categories":[{"name":"test", "test":"test"}],"mechanisms":[]}`).
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
			JSON(`{"name":"test","publisher":"test","playerNumber":1,"tags":[],"categories":[],"mechanisms":[{"name":"test"}]}`).
			Header("Authorization", "Bearer "+suite.base.oauthHeader).
			Expect(suite.T()).
			Status(http.StatusNotFound).
			End()
	apitest.New(). // Invalid Struct -> Mechanisms have too many fields
			HandlerFunc(suite.base.router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"name":"test","publisher":"test","playerNumber":1,"tags":[],"categories":[],"mechanisms":[{"name":"test", "test":"test"}]}`).
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
