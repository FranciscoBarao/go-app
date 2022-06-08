package mechanism

import (
	"catalog/database"
	"catalog/repository/mechanismRepo"
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

	// Set Mechanism Repository
	mechanismRepo := mechanismRepo.NewMechanismRepository(db)
	// Set Mechanism Controller
	controller := InitController(mechanismRepo)

	router = chi.NewRouter()
	router.Post("/api/mechanism", controller.Create)
	router.Get("/api/mechanism", controller.GetAll)
	router.Get("/api/mechanism/{name}", controller.Get)
	router.Delete("/api/mechanism/{name}", controller.Delete)

	log.Println("Setup Complete")
}

func TestPostMechanism(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": "test"}`).
		Expect(t).
		Body(`{"name": "test"}`).
		Status(http.StatusOK).
		End()
}

func TestGetMechanism(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/mechanism/test").
		Expect(t).
		Status(http.StatusOK).
		Body(`{"name": "test"}`).
		End()
}

func TestDeleteMechanism(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/mechanism/test").
		Expect(t).
		Status(http.StatusNoContent).
		End()
}

func TestCreateMechanismFailures(t *testing.T) {
	// Several Json Objects on the body
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`[{"name":"a"},{"name":"b"}]`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{name:"a"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": 1}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"test": "test"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(``).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Invalid Struct -> NOT maxstringlength(30)
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`).
		Expect(t).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> NOT alphanum
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": "test.?"}`).
		Expect(t).
		Status(http.StatusForbidden).
		End()
}

func TestGetMechanismFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/mechanism/test").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestDeleteMechanismFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/mechanism/test").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
