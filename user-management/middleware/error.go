package middleware

import "encoding/json"

type MalformedRequest struct {
	status  int
	message string
}

func NewError(status int, message string) *MalformedRequest {
	return &MalformedRequest{
		status:  status,
		message: message,
	}
}

func (mr *MalformedRequest) Error() string {
	return mr.message
}

func (mr *MalformedRequest) GetStatus() int {
	return mr.status
}

func (mr *MalformedRequest) GetMessage() string {
	b, _ := json.Marshal(mr)
	return string(b)
}