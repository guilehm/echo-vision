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
		Expect(e.Result()).To(Equal(validEvent.Result()))
		Expect(e.Status()).To(Equal(validEvent.Status()))
	})

	It("SaveEvent with file", func() {
		// Arrange
		ve := validEventWithFile(validUser)
		eventID, err := repo.SaveEvent(ctx, nil, ve)
		Expect(err).ToNot(HaveOccurred())

		// Assert
		e, err := repo.FindEventByID(ctx, nil, eventID)
		Expect(err).ToNot(HaveOccurred())

		Expect(e.ID().String()).To(Equal(ve.ID().String()))
		Expect(e.File).ToNot(BeNil())
		Expect(e.File().Filename).To(Equal(ve.File().Filename))
		Expect(e.File().Filepath).To(Equal(ve.File().Filepath))
		Expect(e.File().Filesize).To(Equal(ve.File().Filesize))
		Expect(e.File().ContentType).To(Equal(ve.File().ContentType))
	})
})
