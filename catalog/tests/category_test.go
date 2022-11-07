package tests

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
)

func TestPostCategory(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": "test"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Body(`{"name": "test"}`).
		Status(http.StatusOK).
		End()
}

func TestGetCategory(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/category/test").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		Body(`{"name": "test"}`).
		End()
}

func TestDeleteCategory(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/category/test").
		Header("Authorization", "Bearer "+oauthHeader).
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
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{name:"a"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": 1}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"test": "test"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(``).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Invalid Struct -> NOT maxstringlength(30)
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> NOT alphanum
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": "test.?"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()
}

func TestGetCategoryFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/category/test").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestDeleteCategoryFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/category/test").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
