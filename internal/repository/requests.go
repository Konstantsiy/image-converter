package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Konstantsiy/image-converter/pkg/logger"
)

var ErrNoSuchRequest = errors.New("request with the given id does not exists")

const (
	RequestStatusProcessing = "processing"
	RequestStatusFailed     = "failed"
	RequestStatusDone       = "done"
)

// ConversionRequest represents conversion request in the database.
type ConversionRequest struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	SourceID     string    `json:"source_id"`
	TargetID     string    `json:"target_id"`
	SourceFormat string    `json:"source_format"`
	TargetFormat string    `json:"target_format"`
	Ratio        int       `json:"ratio"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	Status       string    `json:"status"`
}

// RequestsRepository represents repository fro working with requests.
type RequestsRepository struct {
	db *sql.DB
}

func NewRequestsRepository(db *sql.DB) *RequestsRepository {
	return &RequestsRepository{db: db}
}

// GetRequestsByUserID gets the information about requests by given user id.
func (rr *RequestsRepository) GetRequestsByUserID(ctx context.Context, userID string) ([]ConversionRequest, error) {
	var requests []ConversionRequest
	var request ConversionRequest
	var targetIDNull sql.NullString
	const query = `select id, user_id, source_id, target_id, source_format, target_format, ratio, status, created, updated
		from converter.requests where user_id = $1;`

	rows, err := rr.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get user requests: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&request.ID,
			&request.UserID,
			&request.SourceID,
			&targetIDNull,
			&request.SourceFormat,
			&request.TargetFormat,
			&request.Ratio,
			&request.Status,
			&request.Created,
			&request.Updated)
		if err != nil {
			return nil, fmt.Errorf("can't scan user request from rows: %w", err)
		}
		if targetIDNull.Valid {
			request.TargetID = targetIDNull.String
		} else {
			request.TargetID = ""
		}
		requests = append(requests, request)
	}

	if err = rows.Err(); err != nil {
		return requests, fmt.Errorf("error selecting rows: %w", err)
	}

	return requests, nil
}

// InsertRequest creates the conversion request and returns its id.
func (rr *RequestsRepository) InsertRequest(ctx context.Context, userID, sourceID, sourceFormat, targetFormat string, ratio int) (string, error) {
	var requestID string

	const query = `insert into converter.requests 
		(user_id, source_id, target_id, source_format, target_format, ratio, status)
		values ($1, $2, NULL, $3, $4, $5, 'queued') 
		returning id;`

	err := rr.db.QueryRowContext(ctx, query, userID, sourceID, sourceFormat, targetFormat, ratio).Scan(&requestID) // todo
	switch {
	case err == sql.ErrNoRows:
		logger.FromContext(ctx).Infoln("--- no request")
	case err != nil:
		logger.FromContext(ctx).Infoln("--- query error: " + err.Error())
	}
	if err != nil {
		return "", fmt.Errorf("can't make request: %w", err)
	}

	return requestID, nil
}

// UpdateRequest updates the request status and the id of the target image.
func (rr *RequestsRepository) UpdateRequest(ctx context.Context, requestID, status, targetID string) error {
	var sqlTargetID sql.NullString
	if targetID != "" {
		sqlTargetID = sql.NullString{String: targetID, Valid: true}
	}

	const query = "update converter.requests set target_id=$2, status=$3, updated=default where id=$1;"
	res, err := rr.db.ExecContext(ctx, query, requestID, sqlTargetID, status)
	if err != nil {
		return fmt.Errorf("can't update request: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get the number of rows affected by an update: %w", err)
	}
	if count == 0 {
		return ErrNoSuchRequest
	}

	return nil
}
