package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
)

// MalformedRequest represents an HTTP error with a status code and message.
type MalformedRequest struct {
	Status  int
	Message string
}

// NewError creates a new MalformedRequest error.
func NewError(status int, message string) *MalformedRequest {
	return &MalformedRequest{
		Status:  status,
		Message: message,
	}
}

func (mr *MalformedRequest) Error() string {
	return mr.Message
}

// GetStatus returns the HTTP status code.
func (mr *MalformedRequest) GetStatus() int {
	return mr.Status
}

// GetMessage returns the JSON-encoded error message.
func (mr *MalformedRequest) GetMessage() string {
	b, _ := json.Marshal(mr)
	return string(b)
}

// IsNotFound reports whether err is a MalformedRequest with a 404 status.
func IsNotFound(err error) bool {
	if mr, ok := errors.AsType[*MalformedRequest](err); ok {
		return mr.Status == http.StatusNotFound
	}
	return false
}
