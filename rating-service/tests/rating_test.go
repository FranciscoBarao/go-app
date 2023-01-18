package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"rating-service/model"

	"github.com/steinfletcher/apitest"
)

var id string

/* Tests POST a Rating with success*/
func TestCreateRating(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/rating").
		JSON(`{"username":"test", "reference_namespace": "test", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": 10}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Body(`{"username":"test", "reference_namespace": "test", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": 10}`).
		Status(http.StatusOK).
		Assert(func(res *http.Response, req *http.Request) error { // Gets ID from Rating Creation for further test use
			var rating model.Rating
			json.NewDecoder(res.Body).Decode(&rating)
			id = rating.GetId().String()
			return nil
		}).
		End()
}

/* Tests GET all Ratings with success*/
func TestGetAllRatings(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/rating").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		End()
}

/* Tests GET a Rating with success*/
func TestGetRating(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/rating/"+id).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		Body(`{"username":"test", "reference_namespace": "test", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": 10}`).
		End()
}

/* Tests DELETE a Rating with success*/
/* func TestDeleteRating(t *testing.T) {
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/rating/"+id).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusNoContent).
		End()
}*/

/* Tests POST Ratings with Malformed JSON errors*/
func TestCreateRatingFailureMalformedJson(t *testing.T) {
	// Malformed Json - Missing username field
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"reference_namespace": "test", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": 10}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json - Missing reference_namespace field
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"username":"test", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": 10}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json - Missing reference_id field
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"username":"test", "reference_namespace": "test", "value": 10}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json - Missing value field
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"username":"test", "reference_namespace": "test", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()
}

/* Tests POST a Rating with unmarshall type errors*/
func TestCreateRatingFailureUnmarshallType(t *testing.T) {
	// Unmarshall type error - username not string
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"username":123, "reference_namespace": "test", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": 10}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error - reference_namespace not string
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"username":"test", "reference_namespace": 123, "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": 10}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error - reference_id not string
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"username":"test", "reference_namespace": "test", "reference_id": 123, "value": 10}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error - value not int
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"username":"test", "reference_namespace": "test", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": "10"}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()
}

/* Tests POST a Rating with invalid struct errors*/
func TestCreateRatingFailureInvalidStruct(t *testing.T) {
	// Invalid Struct -> username NOT maxstringlength(50)
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"username":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "reference_namespace": "test", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": 10}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> username NOT Alphanumeric
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"username":"asd.?_", "reference_namespace": "test", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": 10}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> reference_namespace NOT maxstringlength(50)
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"username":"test", "reference_namespace": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": 10}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> reference_namespace NOT Alphanumeric
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"username":"test", "reference_namespace": "asd.?_", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": 10}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> reference_id NOT UUID
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"username":"test", "reference_namespace": "test", "reference_id": "test", "value": 10}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> value ABOVE range 0-10
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"username":"test", "reference_namespace": "test", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": 11}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> value BELLOW range 0-10
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/category").
		JSON(`{"username":"test", "reference_namespace": "test", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": -1}`).
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()
}

/* Tests POST a Rating with more generic JSON body related errors*/
func TestCreateRatingJsonFailures(t *testing.T) {
	// Several Json Objects on the body
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Post("/api/rating").
		JSON(`[{"username":"test", "reference_namespace": "test", "reference_id": "f299b8d9-135a-42df-9db8-d1d920d6f456", "value": 10}]`).
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
}

/* Tests GET a Rating with not found error*/
func TestGetRatingFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Get("/api/rating/test").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

/* Tests DELETE a Rating with not found error*/
func TestDeleteRatingFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(router.ServeHTTP).
		Delete("/api/rating/test").
		Header("Authorization", "Bearer "+oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
