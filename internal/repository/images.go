package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ErrNoSuchImage notifies that needed image does not exist.
var ErrNoSuchImage = errors.New("the image with this id does not exist")

// ImagesRepository represents repository fro working with images.
type ImagesRepository struct {
	db *sql.DB
}

// NewImagesRepository creates new images repository.
func NewImagesRepository(db *sql.DB) *ImagesRepository {
	return &ImagesRepository{db: db}
}

// InsertImage inserts the image into images table and returns image id.
func (ir *ImagesRepository) InsertImage(ctx context.Context, filename, format string) (string, error) {
	var imageID string
	const query = "insert into converter.images (name, format) values ($1, $2) returning id;"

	err := ir.db.QueryRowContext(ctx, query, filename, format).Scan(&imageID)
	if err != nil {
		return "", fmt.Errorf("can't insert image: %w", err)
	}

	return imageID, nil
}

// GetImageIDByUserID returns the image id to the storage.
func (ir *ImagesRepository) GetImageIDByUserID(ctx context.Context, userID, imageID string) (string, error) {
	var resImageID string
	const query = `select i.id from converter.requests r
    	join converter.images i
    	on i.id = $2
    	and (r.source_id = i.id or r.target_id = i.id)
    	and r.user_id = $1;`

	err := ir.db.QueryRowContext(ctx, query, userID, imageID).Scan(&resImageID)
	if err == sql.ErrNoRows {
		return "", ErrNoSuchImage
	}
	if err != nil {
		return "", fmt.Errorf("error in the image selection: %w", err)
	}

	return imageID, nil
}
