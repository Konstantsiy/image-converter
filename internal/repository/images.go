package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	// ErrNoSuchImage notifies that needed image does not exist.
	ErrNoSuchImage = errors.New("the image with this id does not exist")
	// ErrEmptySQLDriver notifies that an empty driver is passed to the constructor.
	ErrEmptySQLDriver = errors.New("the SQL-driver for the constructor must not be nil")
)

// ImagesRepository represents repository fro working with images.
type ImagesRepository struct {
	db *sql.DB
}

// NewImagesRepository creates new images repository.
func NewImagesRepository(db *sql.DB) (*ImagesRepository, error) {
	if db == nil {
		return nil, ErrEmptySQLDriver
	}
	return &ImagesRepository{db: db}, nil
}

// InsertImage inserts the image into images table and returns image id.
func (ir *ImagesRepository) InsertImage(ctx context.Context, filename, format string) (string, error) {
	var imageID string
	const query = "INSERT INTO converter.images (name, format) VALUES ($1, $2) RETURNING id;"

	err := ir.db.QueryRowContext(ctx, query, filename, format).Scan(&imageID)
	if err != nil {
		return "", fmt.Errorf("can't insert image: %w", err)
	}

	return imageID, nil
}

// GetImageIDByUserID returns the image id to the storage.
func (ir *ImagesRepository) GetImageIDByUserID(ctx context.Context, userID, imageID string) (string, error) {
	var resImageID string
	const query = `SELECT i.id FROM converter.requests r
    JOIN converter.images i
    ON i.id = $2
    AND (r.source_id = i.id OR r.target_id = i.id)
    AND r.user_id = $1;`

	err := ir.db.QueryRowContext(ctx, query, userID, imageID).Scan(&resImageID)
	if err == sql.ErrNoRows {
		return "", ErrNoSuchImage
	}
	if err != nil {
		return "", fmt.Errorf("error in the image selection: %w", err)
	}

	return imageID, nil
}
