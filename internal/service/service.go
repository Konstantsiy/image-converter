package service

import (
	"github.com/Konstantsiy/image-converter/internal/model"
)

// Users provides interface for UserService.
type Users interface {
	IsExists(email string) bool
	LogIn(request model.AuthRequest) (model.Tokens, error)
	SignUp(request model.AuthRequest) error
}

// Images provides interface for ImageService.
type Images interface {
	Convert(request model.ConversionRequest) (int, error)
	DownloadImage(id string) (string, error)
	GetImageByID(id string) (model.Image, error)
	GetDownloadedURLByLocation(location string) (string, error)
}

// Requests provides interface for RequestService.
type Requests interface {
	GetRequestsHistory(userID string) ([]model.ConversionRequestInfo, error)
}
