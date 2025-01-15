package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/shared"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name        string // description of this test case
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
			id:          uuid.New(),
			firstName:   "John",
			lastName:    "Doe",
			email:       "john.doe@example.com",
			createdAt:   time.Now(),
			updatedAt:   time.Now(),
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "Invalid user ID",
			id:          uuid.Nil,
			firstName:   "Jane",
			lastName:    "Doe",
			email:       "jane.doe@example.com",
			createdAt:   time.Now(),
			updatedAt:   time.Now(),
			wantErr:     true,
			expectedErr: shared.ErrInvalidID,
		},
		{
			name:        "Invalid email format",
			id:          uuid.New(),
			firstName:   "Alice",
			lastName:    "Smith",
			email:       "alice.smith@@example.com",
			createdAt:   time.Now(),
			updatedAt:   time.Now(),
			wantErr:     true,
			expectedErr: shared.ErrInvalidEmail,
		},
		{
			name:        "Empty email",
			id:          uuid.New(),
			firstName:   "Bob",
			lastName:    "Brown",
			email:       "",
			createdAt:   time.Now(),
			updatedAt:   time.Now(),
			wantErr:     true,
			expectedErr: shared.ErrInvalidEmail,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := domain.NewUser(tt.id, tt.firstName, tt.lastName, tt.email, tt.createdAt, tt.updatedAt)
			gotErr := u.Validate()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Validate() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Validate() succeeded unexpectedly")
			}
		})
	}
}
