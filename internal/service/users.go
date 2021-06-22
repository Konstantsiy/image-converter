package service

import (
	"github.com/Konstantsiy/image-converter/internal/model"
)

// UserService provides functionality for working with users.
type UserService struct{}

// LogIn implements user authentication.
func (s *UserService) LogIn(request model.AuthRequest) (model.Tokens, error) {
	// delegating work to the repository
	return model.Tokens{}, nil
}

// SignUp implements user registration.
func (s *UserService) SignUp(request model.AuthRequest) (int, error) {
	// delegating work to the repository
	return -1, nil
}

// IsExists checks if there is a user with a similar email address in the database.
func (s *UserService) IsExists(email string) bool {
	// delegating work to the repository
	return false
}
