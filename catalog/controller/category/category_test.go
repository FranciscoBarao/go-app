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
