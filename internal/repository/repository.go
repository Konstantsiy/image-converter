// Package repository provides the logic for working with database.
package repository

import (
	"errors"
	"time"
)

var (
	ErrNoSuchUser        = errors.New("the user with this email does not exist")
	ErrNoSuchImage       = errors.New("the image with this id does not exist")
	ErrUserAlreadyExists = errors.New("the user with the given username already exists")
)

// Image represents an image in the database.
type Image struct {
	ID       string
	Name     string
	Format   string
	Location string
}

// User represents the user in the database.
type User struct {
	ID       string
	Email    string
	Password string
}

// ConversionRequest represents conversion request in the database.
type ConversionRequest struct {
	ID           string
	Name         string
	UserID       string
	SourceID     string
	TargetID     string
	SourceFormat string
	TargetFormat string
	Ratio        int
	Created      time.Time
	Updated      time.Time
	Status       string
}

// Repository represents the layer between the business logic and the database.
type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

// InsertUser inserts the user into users table.
func (r *Repository) InsertUser(email, password string) (string, error) {
	return "", nil
}

// GetImageLocationByID gets the information about the user by given email.
func (r *Repository) GetUserByEmail(email string) (User, error) {
	return User{}, nil
}

// GetImageLocationByID gets the information about the image by given id.
func (r *Repository) GetImageLocationByID(imageId string) (string, error) {
	return "", nil
}

// GetRequestsByUserID gets the information about requests by given user id.
func (r *Repository) GetRequestsByUserID(userID string) ([]ConversionRequest, error) {
	return []ConversionRequest{}, nil
}
