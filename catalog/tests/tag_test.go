package tests

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
)

func TestCreateTagSuccess(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"name": "test"}`).
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Body(`{"name": "test"}`).
		Status(http.StatusOK).
		End()
}

func TestGetTagSuccess(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/tag/test").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		Body(`{"name": "test"}`).
		End()
}

func TestGetAllTagSuccess(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/tag").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestDeleteTagSuccess(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/tag/test").
		Header("Authorization", "Bearer "+oauthHeader).
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
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(`{name:"a"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"name": 1}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"test": "test"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(``).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Invalid Struct -> NOT maxstringlength(30)
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> NOT alphanum
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"name": "test.?"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()
}

func TestGetTagFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/tag/test").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestDeleteTagFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/tag/test").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
