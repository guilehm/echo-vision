package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	. "github.com/onsi/gomega"
)

func saveUser(u *domain.User) *domain.User {
	hp, _ := passwordAdapter.HashPassword("awesomepassword")
	u.SetHashedPassword(hp)
	vu, err := userUseCase.CreateUser(
		ctx,
		u.FirstName(),
		u.LastName(),
		u.Email(),
		u.Password(),
	)
	Expect(err).ToNot(HaveOccurred())
	_, err = userUseCase.SaveUser(ctx, vu)
	Expect(err).ToNot(HaveOccurred())
	return vu
}

func saveEvent(e *domain.Event) *domain.Event {
	ve, err := eventUseCase.CreateEvent(
		ctx,
		e.UserID(),
		e.EventType().String(),
		e.SubType().String(),
		nil,
	)
	Expect(err).ToNot(HaveOccurred())
	_, err = eventUseCase.SaveEvent(ctx, ve)
	Expect(err).ToNot(HaveOccurred())
	return ve
}

func toReader[T any](t T) io.Reader {
	b, err := json.Marshal(t)
	Expect(err).ToNot(HaveOccurred())
	return bytes.NewReader(b)
}

func generateRefreshTokenMock(user *domain.User, delay time.Duration) string {
	claims := jwt.MapClaims{
		"userID": user.ID(),
		"exp":    time.Now().Add(2 * time.Hour).Add(-delay).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(jwtSecretKey))
	Expect(err).ToNot(HaveOccurred())
	return tokenStr
}
