package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNoSuchImage = errors.New("the image with this id does not exist")

// ImagesRepository represents repository fro working with images.
type ImagesRepository struct {
	db *sql.DB
}

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

// GetImageIDInStore returns the image id to the storage.
func (ir *ImagesRepository) GetImageIDInStore(ctx context.Context, id string) (string, error) {
	var imageID string
	const query = "select id from converter.images where id=$1;"

	err := ir.db.QueryRowContext(ctx, query, id).Scan(&imageID)
	if err == sql.ErrNoRows {
		return "", ErrNoSuchImage
	}
	if err != nil {
		return "", fmt.Errorf("error in the image selection: %w", err)
	}

	return imageID, nil
}
