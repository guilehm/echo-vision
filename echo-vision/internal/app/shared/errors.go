package shared

import "errors"

var (
	ErrInvalidID        = errors.New("invalid id")
	ErrInvalidStatus    = errors.New("invalid status")
	ErrInvalidEventType = errors.New("invalid event type")
	ErrInvalidPayload   = errors.New("invalid payload")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrUserNotLoaded    = errors.New("user not loaded")
)
