package tag

import (
	"catalog/database"
	"catalog/repository/tagRepo"
	"log"
	"net/http"
	"testing"

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

	// Set Tag Repository
	tagRepo := tagRepo.NewTagRepository(db)
	// Set Tag Controller
	controller := InitController(tagRepo)

	router = chi.NewRouter()
	router.Post("/api/tag", controller.Create)
	router.Get("/api/tag", controller.GetAll)
	router.Get("/api/tag/{name}", controller.Get)
	router.Delete("/api/tag/{name}", controller.Delete)

	log.Println("Setup Complete")
}

func TestCreateTagSuccess(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"name": "test"}`).
		Header("Content-Type", "application/json").
		Expect(t).
		Body(`{"name": "test"}`).
		Status(http.StatusOK).
		End()
}

func TestGetTagSuccess(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/tag/test").
		Expect(t).
		Status(http.StatusOK).
		Body(`{"name": "test"}`).
		End()
}

func TestGetAllTagSuccess(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/tag").
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestDeleteTagSuccess(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/tag/test").
		Expect(t).
		Status(http.StatusNoContent).
		End()
}

func TestCreateTagFailures(t *testing.T) {
	// Several Json Objects on the body
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(`[{"name":"a"},{"name":"b"}]`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(`{name:"a"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"name": 1}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"test": "test"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(``).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Invalid Struct -> NOT maxstringlength(30)
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`).
		Expect(t).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> NOT alphanum
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"name": "test.?"}`).
		Expect(t).
		Status(http.StatusForbidden).
		End()
}

func TestGetTagFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/tag/test").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestDeleteTagFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/tag/test").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
