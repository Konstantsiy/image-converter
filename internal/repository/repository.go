// Package repository provides the logic for working with database.
package repository

import "context"

// Users represents users repository.
type Users interface {
	InsertUser(ctx context.Context, email, password string) (string, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

// Images represents images repository.
type Images interface {
	InsertImage(ctx context.Context, filename, format string) (string, error)
	GetImageIDByUserID(ctx context.Context, id string) (string, error)
}

// Requests represents requests repository.
type Requests interface {
	InsertRequest(ctx context.Context, userID, sourceID, sourceFormat, targetFormat string, ratio int) (string, error)
	GetRequestsByUserID(ctx context.Context, userID string) ([]ConversionRequest, error)
	UpdateRequest(ctx context.Context, requestID, status, targetID string) error
}
