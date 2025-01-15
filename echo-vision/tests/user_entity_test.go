package tests

import (
	"time"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/shared"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Validation", func() {
	var (
		validID   = uuid.New()
		validTime = time.Now()
	)

	tests := []struct {
		name        string
		id          uuid.UUID
		firstName   string
		lastName    string
		email       string
		createdAt   time.Time
		updatedAt   time.Time
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "Valid user",
			id:          validID,
			firstName:   "John",
			lastName:    "Doe",
			email:       "john.doe@example.com",
			createdAt:   validTime,
			updatedAt:   validTime,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "Invalid user ID",
			id:          uuid.Nil,
			firstName:   "Jane",
			lastName:    "Doe",
			email:       "jane.doe@example.com",
			createdAt:   validTime,
			updatedAt:   validTime,
			wantErr:     true,
			expectedErr: shared.ErrInvalidID,
		},
		{
			name:        "Invalid email format",
			id:          validID,
			firstName:   "Alice",
			lastName:    "Smith",
			email:       "alice.smith@@example.com",
			createdAt:   validTime,
			updatedAt:   validTime,
			wantErr:     true,
			expectedErr: shared.ErrInvalidEmail,
		},
		{
			name:        "Empty email",
			id:          validID,
			firstName:   "Bob",
			lastName:    "Brown",
			email:       "",
			createdAt:   validTime,
			updatedAt:   validTime,
			wantErr:     true,
			expectedErr: shared.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		It(tt.name, func() {
			u := domain.NewUser(tt.id, tt.firstName, tt.lastName, tt.email, tt.createdAt, tt.updatedAt)
			gotErr := u.Validate()

			if tt.wantErr {
				Expect(gotErr).To(HaveOccurred())
				Expect(gotErr.Error()).To(Equal(tt.expectedErr.Error()))
			} else {
				Expect(gotErr).NotTo(HaveOccurred())
			}
		})
	}
})
