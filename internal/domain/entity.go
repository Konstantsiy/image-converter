// Package domain represents the main entities that are used in the application.
package domain

import "time"

// ImageInfo represents an image in the database.
type ImageInfo struct {
	ID       string
	Name     string
	Format   string
	Location string
}

// UserInfo represents the user in the database.
type UserInfo struct {
	ID       string
	Email    string
	Password string
}

// ConversionRequestInfo represents conversion request in the database.
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
