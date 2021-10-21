// Package service implements a business logic layer.
package service

import (
	"context"
	"mime/multipart"

	"github.com/Konstantsiy/image-converter/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mock/mock.go

// ServiceError represents service related error.
type ServiceError struct {
	Err        error
	StatusCode int
}

func (e *ServiceError) Error() string {
	return e.Err.Error()
}

// Authorization represents authorization service.
type Authorization interface {
	ParseToken(accessToken string) (string, error)
	LogIn(ctx context.Context, email, password string) (string, string, error)
	SignUp(ctx context.Context, email, password string) (string, error)
}

// Images represents images service.
type Images interface {
	Convert(ctx context.Context, sourceFile multipart.File, filename, sourceFormat, targetFormat string, ratio int) (string, string, error)
	Download(ctx context.Context, id string) (string, error)
}

// Requests represents requests service.
type Requests interface {
	GetUsersRequests(ctx context.Context) ([]repository.ConversionRequest, error)
}

// Producer represents queue producer.
type Producer interface {
	SendToQueue(fileID, filename, sourceFormat, targetFormat, requestID string, ratio int) error
}
