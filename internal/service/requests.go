package service

import (
	"github.com/Konstantsiy/image-converter/internal/model"
)

// RequestService provides functionality for working with user requests.
type RequestService struct{}

// GetRequestsHistory returns all user requests by user ID.
func (s *RequestService) GetRequestsHistory(userID string) ([]model.ConversionRequestInfo, error) {
	// delegating work to the repository
	return nil, nil
}
