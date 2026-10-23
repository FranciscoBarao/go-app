package user

// RegisterRequest is the request DTO for POST /api/register.
type RegisterRequest struct {
	Username string `json:"username" valid:"required, alphanum, maxstringlength(50)"`
	Email    string `json:"email" valid:"required, email"`
	Password string `json:"password" valid:"required, minstringlength(6)"`
}

// NewUser creates a User domain model from a RegisterRequest.
func NewUser(req *RegisterRequest) *User {
	return &User{
		Username: req.Username,
		Email:    req.Email,
	}
}

// DeleteRequest is the request DTO for DELETE /api/user/{username}.
type DeleteRequest struct {
	Username string `json:"username" valid:"required, alphanum, maxstringlength(50)"`
}
