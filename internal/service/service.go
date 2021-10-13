package service

import (
	"context"

	"github.com/Konstantsiy/image-converter/internal/repository"
)

type ServiceError struct {
	Err        error
	StatusCode int
}

func (e *ServiceError) Error() string {
	return e.Err.Error()
}

type Authorization interface {
	ParseToken(accessToken string) (string, error)
	LogIn(ctx context.Context, email, password string) (repository.User, error)
	SignUp(ctx context.Context, email, password string) (string, error)
}

type Images interface {
	Convert(ctx context.Context, filename, sourceFormat, targetFormat, ratio string) (string, string, error)
	Download(ctx context.Context, userID, imageID string) (string, error)
}

type Requests interface {
	GetUsersRequests(ctx context.Context) ([]repository.ConversionRequest, error)
}
