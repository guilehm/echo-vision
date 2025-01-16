package tests

import (
	"github.com/guilehm/echo-vision/internal/app/domain"
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
	)
	Expect(err).ToNot(HaveOccurred())
	_, err = eventUseCase.SaveEvent(ctx, ve)
	Expect(err).ToNot(HaveOccurred())
	return ve
}
