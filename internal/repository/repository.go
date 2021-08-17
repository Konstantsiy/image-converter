// Package repository provides the logic for working with database.
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

var (
	ErrNoSuchUser        = errors.New("the user with this email does not exist")
	ErrNoSuchImage       = errors.New("the image with this id does not exist")
	ErrUserAlreadyExists = errors.New("the user with the given email already exists")
	ErrNoSuchRequest     = errors.New("request with the given id does not exists")
)

const uniqueViolationCode = "23505"

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
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// InsertUser inserts the user into users table and returns user id.
func (r *Repository) InsertUser(email, password string) (string, error) {
	var userID string
	const query = "insert into converter.users (email, password) values ($1, $2) returning id;"

	err := r.db.QueryRow(query, email, password).Scan(&userID)
	if err, ok := err.(*pq.Error); ok {
		if err.Code == uniqueViolationCode {
			return "", ErrUserAlreadyExists
		} else {
			return "", fmt.Errorf("can't insert user: %s", err.Code.Name())
		}
	}

	return userID, nil
}

// InsertImage inserts the image into images table and returns image id.
func (r *Repository) InsertImage(filename, format string) error {
	var imageID string
	const query = "insert into converter.images (name, format) values ($1, $2) returning id;"

	err := r.db.QueryRow(query, filename, format).Scan(&imageID)
	if err != nil {
		return fmt.Errorf("can't insert image: %w", err)
	}

	return nil
}

// GetUserByEmail gets the information about the user by given email.
func (r *Repository) GetUserByEmail(email string) (User, error) {
	var user User
	const query = "select id, email, password from converter.users where email = $1;"

	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return User{}, ErrNoSuchUser
	}

	return user, nil
}

// GetImageByID gets the information about the image by given id.
func (r *Repository) GetImageByID(imageId string) (Image, error) {
	var image Image
	const query = "select id, name, format, location from converter.images where id=$1;"

	err := r.db.QueryRow(query, imageId).Scan(&image.ID, &image.Name, &image.Format)
	if err == sql.ErrNoRows {
		return Image{}, ErrNoSuchImage
	}

	return image, nil
}

// GetRequestsByUserID gets the information about requests by given user id.
func (r *Repository) GetRequestsByUserID(userID string) ([]ConversionRequest, error) {
	var requests []ConversionRequest
	var request ConversionRequest
	const query = `select id, name, source_id, target_id, source_format, target_format, ratio, status, created, updated
		from converter.requests where user_id = $1;`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			request.ID,
			request.UserID,
			request.Name,
			request.SourceID,
			request.TargetID,
			request.SourceFormat,
			request.TargetFormat,
			request.Ratio,
			request.Status,
			request.Created,
			request.Updated)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	if err = rows.Err(); err != nil {
		return requests, err
	}

	return requests, nil
}

// MakeRequest creates the conversion request and returns its id.
func (r *Repository) MakeRequest(filename, userID, sourceFormat, targetFormat string, ratio int) (string, error) {
	tx, txErr := r.db.Begin()
	if txErr != nil {
		return "", txErr
	}

	var requestID string
	const (
		insertImageQuery   = "insert into converter.images (name, format) values ($1, $2) returning id;"
		insertRequestQuery = `insert into converter.requests 
		(user_id, source_id, target_id, source_format, target_format, ratio, 'queued')
		values ($1, $2, NULL, $3, $4, $5) 
		returning id;`
	)

	var imageID string
	txErr = tx.QueryRow(insertImageQuery, filename, sourceFormat).Scan(&imageID)
	if txErr != nil {
		tx.Rollback()
		return "", txErr
	}

	txErr = tx.QueryRow(insertRequestQuery, userID, imageID, sourceFormat, targetFormat, ratio).Scan(&requestID)
	if txErr != nil {
		tx.Rollback()
		return "", txErr
	}

	return requestID, tx.Commit()
}

// UpdateRequest updates the request status and the id of the target image.
func (r *Repository) UpdateRequest(requestID, status, targetID string) error {
	var sqlTargetID sql.NullString
	if targetID != "" {
		sqlTargetID = sql.NullString{String: targetID, Valid: true}
	}

	const query = "update converter.requests set target_id=$2, status=$3, updated=default where id=$1;"
	res, err := r.db.Exec(query, requestID, sqlTargetID, status)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNoSuchRequest
	}

	return nil
}
