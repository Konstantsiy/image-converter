package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Konstantsiy/image-converter/internal/appcontext"

	"github.com/Konstantsiy/image-converter/internal/repository"
)

// RequestsService implements logic for working with requests.
type RequestsService struct {
	requestsRepo *repository.RequestsRepository
}

func NewRequestsService(requestsRepo *repository.RequestsRepository) *RequestsService {
	return &RequestsService{requestsRepo: requestsRepo}
}

// GetUsersRequests displays the user's request history.
func (rs *RequestsService) GetUsersRequests(ctx context.Context) ([]repository.ConversionRequest, error) {
	userID, ok := appcontext.UserIDFromContext(ctx)
	if !ok {
		return nil, &ServiceError{
			fmt.Errorf("can't get user id from application context"),
			http.StatusInternalServerError,
		}
	}

	requestsHistory, err := rs.requestsRepo.GetRequestsByUserID(ctx, userID)
	if err != nil {
		return nil, &ServiceError{err, http.StatusInternalServerError}
	}
	return requestsHistory, nil
}
