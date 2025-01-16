package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Event Repository", func() {
	var err error

	BeforeEach(func() {
		_, err = repo.SaveUser(ctx, nil, validUser)
		Expect(err).ToNot(HaveOccurred())
	})

	It("SaveEvent and FindEventByID", func() {
		// Arrange
		eventID, err := repo.SaveEvent(ctx, nil, validEvent)
		Expect(err).ToNot(HaveOccurred())

		// Assert
		e, err := repo.FindEventByID(ctx, nil, eventID)
		Expect(err).ToNot(HaveOccurred())

		Expect(e.ID().String()).To(Equal(validEvent.ID().String()))
		Expect(e.EventType()).To(Equal(validEvent.EventType()))
		Expect(e.SubType()).To(Equal(validEvent.SubType()))
		Expect(e.Payload()).To(Equal(validEvent.Payload()))
		Expect(e.Result()).To(Equal(validEvent.Result()))
		Expect(e.Status()).To(Equal(validEvent.Status()))
	})
})
