package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
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

func handleMessage(topic string, message any) {
	payload, err := json.Marshal(message)
	Expect(err).ToNot(HaveOccurred())
	msg := messaging.Message{
		Topic:   topic,
		Payload: payload,
		Headers: map[string]string{},
	}
	response := adapter.Handle(ctx, msg)
	Expect(response).To(Equal(messaging.Success))
}

func handleDeadLetterMessage(topic string, message any) {
	payload, err := json.Marshal(message)
	Expect(err).ToNot(HaveOccurred())
	msg := messaging.Message{
		Topic:   topic,
		Payload: payload,
		Headers: map[string]string{},
	}
	response := adapter.Handle(ctx, msg)
	Expect(response).To(Equal(messaging.DeadLetter))
}

func jsonEqual(a, b json.RawMessage) bool {
	var o1, o2 any
	if err := json.Unmarshal(a, &o1); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &o2); err != nil {
		return false
	}
	return reflect.DeepEqual(o1, o2)
}

func expectMessageCalled[T any](
	topic string,
	done chan any,
	expectFunc func(msg T),
) {
	handler.Mock.On(
		"Handle",
		mock.Anything,
		mock.MatchedBy(func(msg messaging.Message) bool {
			defer GinkgoRecover()
			Expect(msg.Topic).To(Equal(topic))
			return true
		}),
	).Once().
		Return(messaging.Success).
		Run(func(args mock.Arguments) {
			defer GinkgoRecover()
			Expect(args).To(HaveLen(2))
			Expect(args.Get(0)).ToNot(BeNil())
			Expect(args.Get(1)).ToNot(BeNil())

			message, ok := args.Get(1).(messaging.Message)
			Expect(ok).To(BeTrue())
			Expect(message.Topic).To(Equal(topic))

			var msg T
			err := json.Unmarshal(message.Payload, &msg)
			Expect(err).ToNot(HaveOccurred())
			expectFunc(msg)
			close(done)
		})
}
