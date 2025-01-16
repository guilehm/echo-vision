package domain

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/shared"
)

const minPasswordLength = 5

type User struct {
	id           uuid.UUID
	firstName    string
	lastName     string
	email        string
	password     string
	accessToken  string
	refreshToken string
	createdAt    time.Time
	updatedAt    time.Time
}

func NewUser(
	id uuid.UUID,
	firstName string,
	lastName string,
	email string,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		email:     strings.ToLower(email),
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) FirstName() string {
	return u.firstName
}

func (u *User) LastName() string {
	return u.lastName
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}

func (u *User) AccessToken() string {
	return u.accessToken
}

func (u *User) RefreshToken() string {
	return u.refreshToken
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) SetHashedPassword(password string) {
	u.password = password
}

func (u *User) SetAccessToken(accessToken string) {
	u.accessToken = accessToken
}

func (u *User) SetRefreshToken(refreshToken string) {
	u.refreshToken = refreshToken
}

func (u *User) Validate() error {
	if u.id == uuid.Nil {
		return shared.ErrInvalidID
	}
	_, err := mail.ParseAddress(u.email)
	if err != nil {
		return shared.ErrInvalidEmail
	}
	return nil
}

func (u *User) IsValidPassword(password string) error {
	if len(password) < minPasswordLength {
		return shared.ErrInvalidPassword
	}
	return nil
}
