package tests

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
)

func TestPostMechanism(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": "test"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Body(`{"name": "test"}`).
		Status(http.StatusOK).
		End()
}

func TestGetMechanism(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/mechanism/test").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		Body(`{"name": "test"}`).
		End()
}

func TestDeleteMechanism(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/mechanism/test").
		Header("Authorization", "Bearer "+oauthHeader).
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
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{name:"a"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": 1}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"test": "test"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(``).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Invalid Struct -> NOT maxstringlength(30)
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> NOT alphanum
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": "test.?"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()
}

func TestGetMechanismFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/mechanism/test").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestDeleteMechanismFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/mechanism/test").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
