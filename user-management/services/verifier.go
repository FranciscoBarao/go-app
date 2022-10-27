package services

import (
	"errors"
	"net/http"
	"user-management/repositories"

	"github.com/go-chi/oauth"
)

type VerifierService struct {
	userRepo userRepository
}

func InitVerifierService(userRepository *repositories.UserRepository) VerifierService {
	return VerifierService{
		userRepo: userRepository,
	}
}

// ValidateUser validates username and password returning an error if the user credentials are wrong
func (*VerifierService) ValidateUser(username, password, scope string, r *http.Request) error {
	if username == "user01" && password == "12345" {
		return nil
	}

	// Get User from DB. If fetched, then it has registered and we can return the token

	return errors.New("wrong user")
}

// ValidateClient validates clientID and secret returning an error if the client credentials are wrong
func (*VerifierService) ValidateClient(clientID, clientSecret, scope string, r *http.Request) error {
	if clientID == "abcdef" && clientSecret == "12345" {
		return nil
	}

	return errors.New("wrong client")
}

// ValidateCode validates token ID
func (*VerifierService) ValidateCode(clientID, clientSecret, code, redirectURI string, r *http.Request) (string, error) {
	return "", nil
}

// AddClaims provides additional claims to the token
func (*VerifierService) AddClaims(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	claims := make(map[string]string)
	claims["customer_id"] = "1001"
	claims["customer_data"] = `{"order_date":"2016-12-14","order_id":"9999"}`
	return claims, nil
}

// AddProperties provides additional information to the token response
func (*VerifierService) AddProperties(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	props := make(map[string]string)
	props["customer_name"] = "Gopher"
	return props, nil
}

// ValidateTokenID validates token ID
func (*VerifierService) ValidateTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}

// StoreTokenID saves the token id generated for the user
func (*VerifierService) StoreTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}
