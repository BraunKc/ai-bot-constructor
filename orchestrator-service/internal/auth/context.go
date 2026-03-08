package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired     = errors.New("token expired")
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidClaims    = errors.New("invalid claims")
	ErrUserIDIsRequired = errors.New("user_id is required")
)

type TokenManager struct {
	secretKey []byte
}

func NewTokenManager(secretKey string) (*TokenManager, error) {
	if len(secretKey) < 32 {
		return nil, errors.New("secret key must be at least 32 bytes")
	}

	return &TokenManager{
		secretKey: []byte(secretKey),
	}, nil
}

func (tm *TokenManager) VerifyToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (any, error) {
		return tm.secretKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrTokenExpired
		}

		return "", ErrInvalidToken
	}
	if !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidClaims
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", ErrUserIDIsRequired
	}

	return userID, nil
}
