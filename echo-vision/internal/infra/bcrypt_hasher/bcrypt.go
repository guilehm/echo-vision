package bcrypthasher

import (
	"github.com/guilehm/echo-vision/internal/app/ports"
	"github.com/guilehm/echo-vision/internal/app/shared"
	"golang.org/x/crypto/bcrypt"
)

type BcryptAdapter struct{}

func (b *BcryptAdapter) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func (b *BcryptAdapter) ValidatePassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return shared.ErrInvalidPassword
	}
	return nil
}

func NewBcryptAdapter() ports.PasswordManager {
	return &BcryptAdapter{}
}
