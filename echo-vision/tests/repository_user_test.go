package tests

import (
	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/shared"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Repository", func() {
	It("SaveUser", func() {
		// Act
		userID, err := repo.SaveUser(ctx, nil, validUser)
		Expect(err).ToNot(HaveOccurred())

		// Assert
		u, err := repo.FindUserByID(ctx, nil, userID)
		Expect(err).ToNot(HaveOccurred())

		Expect(u.ID().String()).To(Equal(validUser.ID().String()))
		Expect(u.FirstName()).To(Equal(validUser.FirstName()))
		Expect(u.LastName()).To(Equal(validUser.LastName()))
		Expect(u.Email()).To(Equal(validUser.Email()))
	})

	Context("with user saved", func() {
		BeforeEach(func() {
			_, err := repo.SaveUser(ctx, nil, validUser)
			Expect(err).ToNot(HaveOccurred())
		})
		It("FindUserByID", func() {
			u, err := repo.FindUserByID(ctx, nil, validUser.ID())
			Expect(err).ToNot(HaveOccurred())

			Expect(u.ID().String()).To(Equal(validUser.ID().String()))
		})

		It("FindUserByEmail", func() {
			u, err := repo.FindUserByEmail(ctx, nil, validUser.Email())
			Expect(err).ToNot(HaveOccurred())

			Expect(u.ID().String()).To(Equal(validUser.ID().String()))
		})
		It("Returns error when user not found", func() {
			_, err := repo.FindUserByID(ctx, nil, uuid.Nil)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(shared.ErrNotNound))
		})
	})
})
