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

func TestPostTag(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"name": "test"}`).
		Expect(t).
		Body(`{"name": "test"}`).
		Status(http.StatusOK).
		End()
}

func TestGetTag(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/tag/test").
		Expect(t).
		Status(http.StatusOK).
		Body(`{"name": "test"}`).
		End()
}

func TestDeleteTag(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/tag/test").
		Expect(t).
		Status(http.StatusNoContent).
		End()
}
