package service

import (
	"github.com/Konstantsiy/image-converter/internal/model"
)

// ImageService provides functionality for working with images.
type ImageService struct{}

// Convert converts the image and returns the conversion request ID.
func (s *ImageService) Convert(request model.ConversionRequest) (int, error) {
	// delegating work to the repository
	return -1, nil
}

// DownloadImage implements downloading image by its ID.
func (s *ImageService) DownloadImage(id string) (string, error) {
	// delegating work to the repository
	return "", nil
}

// GetImageByID returns an image from the database by its id.
func (s *ImageService) GetImageByID(id string) (model.Image, error) {
	// delegating work to the repository
	return model.Image{}, nil
}

// GetDownloadedURLByLocation returns an image from the storage by its location.
func (s *ImageService) GetDownloadedURLByLocation(location string) (string, error) {
	// delegating work to the storage
	return "", nil
}
