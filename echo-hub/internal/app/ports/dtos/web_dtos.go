package dtos

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventResponse struct {
	UserID    uuid.UUID       `json:"userID"`
	ID        uuid.UUID       `json:"id"`
	EventType string          `json:"eventType"`
	SubType   string          `json:"subType"`
	Status    string          `json:"status"`
	Result    json.RawMessage `json:"result,omitempty"`
	File      *FileResponse   `json:"file,omitempty"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdateAt  time.Time       `json:"updatedAt"`
}

type FileResponse struct {
	URL         string `json:"url"`
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Filesize    int64  `json:"filesize"`
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
