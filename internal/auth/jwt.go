// Package auth provides the logic for working with JWT tokens.
package auth

// TokenManager implements functionality for Access & Refresh tokens generation.
type TokenManager struct {
	publicKey string
}

func NewTokenManager(publicKey string) *TokenManager {
	return &TokenManager{
		publicKey: publicKey,
	}
}

// GenerateAccessToken generates new access token.
func (tm *TokenManager) GenerateAccessToken(userID string) (string, error) {
	return "", nil
}

// GenerateRefreshToken generates new refresh token.
func (tm *TokenManager) GenerateRefreshToken() (string, error) {
	return "", nil
}
