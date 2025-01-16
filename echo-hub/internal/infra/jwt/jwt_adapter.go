package jwtadapter

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/shared"
)

type JWTManager struct {
	secretKey       string
	accessDuration  time.Duration
	refreshDuration time.Duration
}

// NewJWTManager creates a new instance of JWTManager
func NewJWTManager(secretKey string, accessDuration, refreshDuration time.Duration) ports.TokenManager {
	return &JWTManager{
		secretKey:       secretKey,
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

// GenerateAccessToken generates a short-lived Access Token
func (j *JWTManager) GenerateAccessToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"userID":    user.ID(),
		"email":     user.Email(),
		"firstName": user.FirstName(),
		"lastName":  user.LastName(),
		"exp":       time.Now().Add(j.accessDuration).Unix(),
		"nonce":     uuid.NewString(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// GenerateRefreshToken generates a long-lived Refresh Token
func (j *JWTManager) GenerateRefreshToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"userID": user.ID(),
		"exp":    time.Now().Add(j.refreshDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken validates an Access or Refresh Token and returns claims if valid
func (j *JWTManager) ValidateToken(tokenString string) (ports.TokenClaims, error) {
	secret := j.secretKey

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, shared.ErrInvalidSigningMethod
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenClaims := make(ports.TokenClaims)
		for key, value := range claims {
			tokenClaims[key] = value
		}
		return tokenClaims, nil
	}
	return nil, shared.ErrInvalidToken
}

// RefreshAccessToken generates a new Access Token using a valid Refresh Token
func (j *JWTManager) RefreshAccessToken(user *domain.User, refreshToken string) (string, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return "", shared.ErrInvalidRefreshToken
	}

	userID, ok := claims["userID"].(uuid.UUID)
	if !ok {
		return "", shared.ErrInvalidRefreshToken
	}
	if userID.String() != user.ID().String() {
		return "", shared.ErrInvalidRefreshToken
	}
	return j.GenerateAccessToken(user)
}
