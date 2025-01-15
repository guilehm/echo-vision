package domain

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/shared"
)

type User struct {
	id        uuid.UUID
	firstName string
	lastName  string
	email     string
	createdAt time.Time
	updatedAt time.Time
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
		email:     email,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
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
