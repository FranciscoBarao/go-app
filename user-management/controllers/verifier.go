package controllers

import (
	"errors"
	"log"
	"net/http"
	"user-management/services"

	"github.com/go-chi/oauth"
)

type VerifierController struct {
	userService userService
}

func InitVerifierController(userSvc *services.UserService) VerifierController {
	return VerifierController{
		userService: userSvc,
	}
}

// ValidateUser validates username and password returning an error if the user credentials are wrong
func (service *VerifierController) ValidateUser(username, password, scope string, r *http.Request) error {

	return service.userService.Login(username, password)
}

// ValidateClient validates clientID and secret returning an error if the client credentials are wrong
func (*VerifierController) ValidateClient(clientID, clientSecret, scope string, r *http.Request) error {
	if clientID == "id" && clientSecret == "secret" {
		return nil
	}

	return errors.New("wrong client")
}

// ValidateCode validates token ID
func (*VerifierController) ValidateCode(clientID, clientSecret, code, redirectURI string, r *http.Request) (string, error) {
	return "WHAT AM I DOING", nil
}

// AddClaims provides additional claims to the token
func (*VerifierController) AddClaims(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	claims := make(map[string]string)
	claims["username"] = credential
	return claims, nil
}

// AddProperties provides additional information to the token response
func (*VerifierController) AddProperties(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	props := make(map[string]string)
	return props, nil
}

// ValidateTokenID validates token ID
func (*VerifierController) ValidateTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	log.Println(credential)
	return nil
}

// StoreTokenID saves the token id generated for the user
func (*VerifierController) StoreTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}
