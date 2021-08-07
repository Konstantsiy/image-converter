package domain

// AuthRequest represents the user's authorization request.
type AuthRequest struct {
	Email    string
	Password string
}

// ConversionRequest represents an image conversion request.
type ConversionRequest struct {
	File         string
	SourceFormat string
	TargetFormat string
	Ratio        int
}
