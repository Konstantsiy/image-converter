// Package model represents the main entities that are used in the application.
package model

import "time"

// Tokens represents tokens for authorization of requests.
type Tokens struct {
	AccessToken  string
	RefreshToken string
}

// AuthRequest represents the user's authorization request.
type AuthRequest struct {
	Email    string
	Password string
}

// Image represents an image in the database.
type Image struct {
	Name     string
	Format   string
	Location string
}

// ConversionRequest represents an image conversion request.
type ConversionRequest struct {
	File         string
	SourceFormat string
	TargetFormat string
	Ratio        int
}

// ConversionRequestInfo represents information about an image conversion request from database.
type ConversionRequestInfo struct {
	ID           string
	Name         string
	SourceID     string
	TargetID     string
	SourceFormat string
	TargetFormat string
	Ratio        int
	Created      time.Time
	Updated      time.Time
	Status       string
}
