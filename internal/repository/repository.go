// Package repository provides the logic for working with database.
package repository

import "context"

type Users interface {
	InsertUser(ctx context.Context, email, password string) (string, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

type Images interface {
	InsertImage(ctx context.Context, filename, format string) (string, error)
	GetImageIDInStore(ctx context.Context, id string) (string, error)
}

type Requests interface {
	GetRequestsByUserID(ctx context.Context, userID string) ([]ConversionRequest, error)
	InsertRequest(ctx context.Context, userID, sourceID, sourceFormat, targetFormat string, ratio int) (string, error)
	UpdateRequest(ctx context.Context, requestID, status, targetID string) error
}
