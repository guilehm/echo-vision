package shared

import "errors"

var (
	ErrInvalidID           = errors.New("invalid id")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrInvalidEventType    = errors.New("invalid event type")
	ErrInvalidEventPayload = errors.New("invalid event payload")
	ErrInvalidPayload      = errors.New("invalid payload")
	ErrInvalidEmail        = errors.New("invalid email")
	ErrUserNotLoaded       = errors.New("user not loaded")
	ErrUserNotFound        = errors.New("user not found")
	ErrNotFound            = errors.New("not found")
	ErrInvalidRequest      = errors.New("invalid request")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidType         = errors.New("invalid type")
	ErrInvalidSubType      = errors.New("invalid sub type")
	ErrInvalidFile         = errors.New("invalid file")

	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrContextValueNotFound = errors.New("context value not found")
	ErrDecodingRequestBody  = errors.New("decoding request body")

	ErrInvalidCursor     = errors.New("invalid cursor")
	ErrInvalidQueryParam = errors.New("invalid query parameter")
)
