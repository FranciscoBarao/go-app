package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/oauth"
)

// LoginService defines the interface needed by the verifier to validate credentials.
type LoginService interface {
	Login(ctx context.Context, username, password string) error
}

// Verifier implements the go-chi/oauth credential verifier interface.
type Verifier struct {
	loginSvc LoginService
}

// NewVerifier creates a new OAuth Verifier.
func NewVerifier(svc LoginService) *Verifier {
	return &Verifier{loginSvc: svc}
}

// ValidateUser validates username and password.
func (v *Verifier) ValidateUser(username, password, scope string, r *http.Request) error {
	return v.loginSvc.Login(context.Background(), username, password)
}

// ValidateClient is hardcoded and a placeholder to validate client credentials
// TODO: add client credentials validations
func (v *Verifier) ValidateClient(clientID, clientSecret, scope string, r *http.Request) error {
	if clientID == "id" && clientSecret == "secret" {
		return nil
	}
	return errors.New("wrong client")
}

// ValidateCode is a NoOp placeholder method which will validate token ID in the future.
// TODO: add client token ID validation
func (v *Verifier) ValidateCode(clientID, clientSecret, code, redirectURI string, r *http.Request) (string, error) {
	return "", nil
}

// AddClaims provides additional claims to the token.
func (v *Verifier) AddClaims(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	claims := make(map[string]string)
	claims["username"] = credential
	return claims, nil
}

// AddProperties provides additional information to the token response.
func (v *Verifier) AddProperties(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	return make(map[string]string), nil
}

// ValidateTokenID validates token ID.
func (v *Verifier) ValidateTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}

// StoreTokenID saves the token id generated for the user.
func (v *Verifier) StoreTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}
