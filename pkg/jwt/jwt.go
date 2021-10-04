// Package auth provides the logic for working with JWT tokens.
package jwt

import (
	"context"
	"fmt"
	"io/ioutil"
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
	publicKey  []byte
	privateKey []byte
}

func NewTokenManager(conf *config.JWTConfig) (*TokenManager, error) {
	privateKey, err := ioutil.ReadFile(conf.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("can't read private key form file: %w", err)
	}

	publicKey, err := ioutil.ReadFile(conf.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("can't read public key form file: %w", err)
	}

	return &TokenManager{publicKey: publicKey, privateKey: privateKey}, nil
}

// GenerateAccessToken generates new access token.
func (tm *TokenManager) GenerateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(AccessTokenTimeout).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(tm.privateKey)
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}

// GenerateRefreshToken generates new refresh token.
func (tm *TokenManager) GenerateRefreshToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RefreshTokenTimeout).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(tm.privateKey)
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}

//ParseToken parses the given token, checks it validity and returns the user ID.
func (tm *TokenManager) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwt.ParseRSAPublicKeyFromPEM(tm.publicKey)
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("can't parse token: %w", err)
	}

	logger.Info(context.Background(), "token is valid")

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("can't get user claims from token")
	}

	logger.Info(context.Background(), "token -> userID: "+claims.Subject)

	return claims.Subject, nil
}
