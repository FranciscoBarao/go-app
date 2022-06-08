package boardgame

import (
	"log"
	"net/http"
	"testing"

	"catalog/database"
	"catalog/repository/boardgameRepo"
	"catalog/repository/categoryRepo"
	"catalog/repository/mechanismRepo"
	"catalog/repository/tagRepo"

	"github.com/go-chi/chi/v5"
	"github.com/steinfletcher/apitest"
)

var router *chi.Mux

func init() {
	log.Println("Setup Starting")

	// Set Database
	db, err := database.Connect()
	if err != nil {
		log.Println("Error occurred while connecting to database")
		return
	}

	// Set Repositories
	boardgameRepo := boardgameRepo.NewBoardGameRepository(db)
	tagRepo := tagRepo.NewTagRepository(db)
	categoryRepo := categoryRepo.NewCategoryRepository(db)
	mechanismRepo := mechanismRepo.NewMechanismRepository(db)

	// Set Boardgame Controller
	controller := InitController(boardgameRepo, tagRepo, categoryRepo, mechanismRepo)

	router = chi.NewRouter()
	router.Post("/api/boardgame", controller.Create)
	router.Get("/api/boardgame", controller.GetAll)
	router.Get("/api/boardgame/{id}", controller.Get)
	router.Patch("/api/boardgame/{id}", controller.Update)
	router.Delete("/api/boardgame/{id}", controller.Delete)
	router.Post("/api/boardgame/{id}/expansion", controller.Create)

	log.Println("Setup Complete")
}

func TestCreateBoardgameSuccess(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestGetBoardgameSuccess(t *testing.T) {
	// Verifies if the created boardgame exists by ID (Runnable only once)
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/boardgame/29").
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestGetAllBoardgameSuccess(t *testing.T) {
	// Verifies if the created boardgame exists on GetAll
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/boardgame").
		Expect(t).
		Status(http.StatusOK).
		End()
}

/*
func TestDeleteBoardgameSuccess(t *testing.T) {
	// Deletes boardgame by ID (Runnable only once)
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/boardgame/48").
		Expect(t).
		Status(http.StatusNoContent).
		End()
}
*/
func TestCreateBoardgameJsonFailures(t *testing.T) {
	// Several Json Objects on the body
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`[{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]},{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}]`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"Name:"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"Name":100,"Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"TEST":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/boardgame").
		JSON(``).
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
			Expect(t).
			Status(http.StatusForbidden).
			End()

	apitest.New(). // Invalid Struct -> NOT alphanum
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test?","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	//  <<<< field - Publisher >>>>
	apitest.New(). // Invalid Struct -> NOT maxstringlength(100)
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	apitest.New(). // Invalid Struct -> NOT alphanum
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test?","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	//  <<<< field - Price >>>>
	apitest.New(). // Invalid Struct -> NOT float
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":"string","PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

	apitest.New(). // Invalid Struct -> NOT in range(0|1000)
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":1001,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	//  <<<< field - PlayerNumber >>>>
	apitest.New(). // Invalid Struct -> NOT int
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1.5,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

	apitest.New(). // Invalid Struct -> NOT in range(0|16)
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":17,"Tags":[],"Categories":[],"Mechanisms":[]}`).
			Expect(t).
			Status(http.StatusForbidden).
			End()

	//  <<<< field - Tags >>>>
	apitest.New(). // Invalid Struct -> Tag does not previously exist
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[{"name":"test"}],"Categories":[],"Mechanisms":[]}`).
			Expect(t).
			Status(http.StatusNotFound).
			End()

	apitest.New(). // Invalid Struct -> Tags have too many fields
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[{"name":"test", "test":"test"}],"Categories":[],"Mechanisms":[]}`).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

	//  <<<< field - Categories >>>>
	apitest.New(). // Invalid Struct -> Category does not previously exist
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[{"name":"test"}],"Mechanisms":[]}`).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	apitest.New(). // Invalid Struct -> Categories have too many fields
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[{"name":"test", "test":"test"}],"Mechanisms":[]}`).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

	//  <<<< field - Mechanisms >>>>
	apitest.New(). // Invalid Struct -> Mechanism does not previously exist
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[{"name":"test"}]}`).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	apitest.New(). // Invalid Struct -> Mchanisms have too many fields
			HandlerFunc(router.ServeHTTP).
			Post("/api/boardgame").
			JSON(`{"Name":"test","Publisher":"test","Price":10,"PlayerNumber":1,"Tags":[],"Categories":[],"Mechanisms":[{"name":"test", "test":"test"}]}`).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
}

func TestGetBoardgameFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/boardgame/1000").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
func TestDeleteBoardgameFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/boardgame/1000").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
