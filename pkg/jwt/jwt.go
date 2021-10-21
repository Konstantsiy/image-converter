// Package auth provides the logic for working with JWT tokens.
package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/Konstantsiy/image-converter/pkg/logger"

	"github.com/Konstantsiy/image-converter/internal/config"

	"github.com/dgrijalva/jwt-go"
)

const (
	AccessTokenTimeout  = 12 * time.Hour
	RefreshTokenTimeout = 48 * time.Hour
)

// TokenManager implements functionality for Access & Refresh tokens generation.
type TokenManager struct {
	signingKey string
}

func NewTokenManager(conf *config.JWTConfig) (*TokenManager, error) {
	if conf.SigningKey == "" {
		return nil, fmt.Errorf("JWT configuration should not be empty")
	}

	return &TokenManager{signingKey: conf.SigningKey}, nil
}

// GenerateAccessToken generates new access token.
func (tm *TokenManager) GenerateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(AccessTokenTimeout).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	return token.SignedString([]byte(tm.signingKey))
}

// GenerateRefreshToken generates new refresh token.
func (tm *TokenManager) GenerateRefreshToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RefreshTokenTimeout).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	return token.SignedString([]byte(tm.signingKey))
}

//ParseToken parses the given token, checks it validity and returns the user ID.
func (tm *TokenManager) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(tm.signingKey), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("can't parse token: %w", err)
	}

	logger.FromContext(context.Background()).Infoln("token is valid")

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("can't get user claims from token")
	}

	return claims.Subject, nil
}
