// Package repository provides the logic for working with database.
package repository

import (
	"errors"

	"github.com/Konstantsiy/image-converter/internal/domain"
)

var (
	ErrNoSuchUser  = errors.New("the user with this email does not exist")
	ErrNoSuchImage = errors.New("the image with this id does not exists")
)

// Repository represents the layer between the business logic and the database.
type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

// InsertUser inserts the user into users table.
func (r *Repository) InsertUser(email, password string) (string, error) {
	return "", nil
}

// GetImageByID gets the information about the user by given email.
func (r *Repository) GetUserByEmail(email string) (domain.UserInfo, error) {
	return domain.UserInfo{}, nil
}

// GetImageByID gets the information about the image by given id.
func (r *Repository) GetImageByID(imageId string) (domain.ImageInfo, error) {
	return domain.ImageInfo{}, nil
}

// GetRequestsByUserID gets the information about requests by given user id.
func (r *Repository) GetRequestsByUserID(userID string) ([]domain.ConversionRequestInfo, error) {
	return []domain.ConversionRequestInfo{}, nil
}
