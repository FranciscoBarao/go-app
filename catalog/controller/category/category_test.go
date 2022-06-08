package category

import (
	"catalog/database"
	"catalog/repository/categoryRepo"
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

	// Set Category Repository
	categoryRepo := categoryRepo.NewCategoryRepository(db)
	// Set Category Controller
	controller := InitController(categoryRepo)

	router = chi.NewRouter()
	router.Post("/api/category", controller.Create)
	router.Get("/api/category", controller.GetAll)
	router.Get("/api/category/{name}", controller.Get)
	router.Delete("/api/category/{name}", controller.Delete)

	log.Println("Setup Complete")
}

func TestPostCategory(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": "test"}`).
		Expect(t).
		Body(`{"name": "test"}`).
		Status(http.StatusOK).
		End()
}

func TestGetCategory(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/category/test").
		Expect(t).
		Status(http.StatusOK).
		Body(`{"name": "test"}`).
		End()
}

func TestDeleteCategory(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/category/test").
		Expect(t).
		Status(http.StatusNoContent).
		End()
}

func TestCreateCategoryFailures(t *testing.T) {
	// Several Json Objects on the body
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`[{"name":"a"},{"name":"b"}]`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{name:"a"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": 1}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"test": "test"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(``).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Invalid Struct -> NOT maxstringlength(30)
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`).
		Expect(t).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> NOT alphanum
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": "test.?"}`).
		Expect(t).
		Status(http.StatusForbidden).
		End()
}

func TestGetCategoryFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/category/test").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestDeleteCategoryFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/category/test").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
