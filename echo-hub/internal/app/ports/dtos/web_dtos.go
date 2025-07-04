package dtos

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventResponse struct {
	UserID    uuid.UUID       `json:"userId"`
	ID        uuid.UUID       `json:"id"`
	EventType string          `json:"eventType"`
	SubType   string          `json:"subType"`
	Status    string          `json:"status"`
	Result    json.RawMessage `json:"result,omitempty"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdateAt  time.Time       `json:"updatedAt"`
}

type EventCreateResponse struct {
	ID uuid.UUID `json:"id"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
}

type UserCreateResponse struct {
	ID           uuid.UUID `json:"id"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
}

type UserLoginResponse struct {
	ID           uuid.UUID `json:"id"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
}
