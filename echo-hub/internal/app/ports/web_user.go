package ports

import (
	"net/http"
)

type UserWebPort interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	MeUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

type UserCreateInput struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type UserLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRefreshTokenInput struct {
	RefreshToken string `json:"refreshToken"`
}
