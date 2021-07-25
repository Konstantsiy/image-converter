// Package auth provides the logic for working with JWT tokens.
package auth


import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	ACTimeout = 12 * time.Hour
	RTTimeout = 48 * time.Hour
)

// TokenManager implements functionality for Access & Refresh tokens generation.
type TokenManager struct {
	privateKey string
}

//NewTokenManager returns new token manager with the given private key.
func NewTokenManager(privateKey string) *TokenManager {
	return &TokenManager{
		privateKey: privateKey,
	}
}

// GenerateAccessToken generates new access token.
func (tm *TokenManager) GenerateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(ACTimeout).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	return token.SignedString(tm.privateKey)
}

// GenerateRefreshToken generates new refresh token.
func (tm *TokenManager) GenerateRefreshToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RTTimeout).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	return token.SignedString(tm.privateKey)
}

//ParseToken parses the given token and returns the user ID.
func (tm *TokenManager) ParseToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(tm.privateKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}
