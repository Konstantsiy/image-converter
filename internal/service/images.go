package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/Konstantsiy/image-converter/internal/appcontext"
	"github.com/Konstantsiy/image-converter/pkg/logger"

	"github.com/Konstantsiy/image-converter/internal/queue"
	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/internal/storage"
)

type ImageService struct {
	imagesRepo   *repository.ImagesRepository
	requestsRepo *repository.RequestsRepository
	storage      *storage.Storage
	producer     *queue.Producer
}

func NewImageService(imagesRepo *repository.ImagesRepository, requestsRepo *repository.RequestsRepository, storage *storage.Storage, producer *queue.Producer) *ImageService {
	return &ImageService{imagesRepo: imagesRepo, requestsRepo: requestsRepo, storage: storage, producer: producer}
}

// Convert converts needed image according to the request.
func (is *ImageService) Convert(ctx context.Context, sourceFile multipart.File, filename, sourceFormat, targetFormat string, ratio int) (string, string, error) {
	userID, ok := appcontext.UserIDFromContext(ctx)
	if !ok {
		return "", "", fmt.Errorf("can't get user id from application context")
	}

	sourceFileID, err := is.imagesRepo.InsertImage(ctx, filename, sourceFormat)
	if err != nil {
		return "", "", &ServiceError{
			fmt.Errorf("repository error: %w", err),
			http.StatusInternalServerError}
	}
	logger.FromContext(ctx).WithField("file_id", sourceFileID).
		Infoln("original file successfully saved in the database")

	err = is.storage.UploadFile(sourceFile, sourceFileID)
	if err != nil {
		return "", "", &ServiceError{
			fmt.Errorf("storage error: %w", err),
			http.StatusInternalServerError}
	}
	logger.FromContext(ctx).WithField("file_id", sourceFileID).
		Infoln("original file successfully uploaded to the S3 storage")

	requestID, err := is.requestsRepo.InsertRequest(ctx, userID, sourceFileID, sourceFormat, targetFormat, ratio)
	if err != nil {
		return "", "", &ServiceError{
			fmt.Errorf("repository error: %w", err),
			http.StatusInternalServerError}
	}
	logger.FromContext(ctx).WithField("request_id", requestID).
		Infoln("request created with the status \"queued\"")

	return sourceFileID, requestID, nil
}

// Download allows you to download original/converted image by id.
func (is *ImageService) Download(ctx context.Context, id string) (string, error) {
	userID, ok := appcontext.UserIDFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("can't get user id from application context")
	}

	imageID, err := is.imagesRepo.GetImageIDByUserID(ctx, userID, id)
	if errors.Is(err, repository.ErrNoSuchImage) {
		return "", &ServiceError{err, http.StatusNotFound}
	}
	if err != nil {
		return "", &ServiceError{err, http.StatusInternalServerError}
	}

	url, err := is.storage.GetDownloadURL(imageID)
	if err != nil {
		return "", &ServiceError{err, http.StatusInternalServerError}
	}

	return url, nil
}
