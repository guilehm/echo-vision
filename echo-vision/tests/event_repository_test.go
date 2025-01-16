package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Event Repository", func() {
	Context("EventCreation", func() {
		It("Valid event", func() {
			// Arrange
			tx, err := repo.BeginTx(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Act
			err = repo.SaveUser(ctx, tx, validUser)
			Expect(err).ToNot(HaveOccurred())
			eventID, err := repo.SaveEvent(ctx, tx, validEvent)
			Expect(err).ToNot(HaveOccurred())

			// Assert
			e, err := repo.FindEventByID(ctx, tx, eventID)
			Expect(err).ToNot(HaveOccurred())

			Expect(e.ID().String()).To(Equal(validEvent.ID().String()))
			Expect(e.EventType()).To(Equal(validEvent.EventType()))
			Expect(e.SubType()).To(Equal(validEvent.SubType()))
			Expect(e.Payload()).To(Equal(validEvent.Payload()))
			Expect(e.Result()).To(Equal(validEvent.Result()))
			Expect(e.Status()).To(Equal(validEvent.Status()))
		})
	})
})
