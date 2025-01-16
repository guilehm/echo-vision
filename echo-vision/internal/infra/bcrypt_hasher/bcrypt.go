package bcrypthasher

import (
	"github.com/guilehm/echo-vision/internal/app/ports"
	"golang.org/x/crypto/bcrypt"
)

type BcryptAdapter struct{}

func (b *BcryptAdapter) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func (b *BcryptAdapter) ValidatePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func NewBcryptAdapter() ports.PasswordManager {
	return &BcryptAdapter{}
}
