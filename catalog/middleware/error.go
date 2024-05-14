package middleware

import "encoding/json"

type MalformedRequest struct {
	Status  int
	Message string
}

func NewError(status int, message string) *MalformedRequest {
	return &MalformedRequest{
		Status:  status,
		Message: message,
	}
}

func (mr *MalformedRequest) Error() string {
	return mr.Message
}

func (mr *MalformedRequest) GetStatus() int {
	return mr.Status
}

func (mr *MalformedRequest) GetMessage() string {
	b, _ := json.Marshal(mr)
	return string(b)
}
