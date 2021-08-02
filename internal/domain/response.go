package domain

// TokenResponse represents token for authorization response.
type TokensResponse struct {
	AccessToken  string
	RefreshToken string
}
