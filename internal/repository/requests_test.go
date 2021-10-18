package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRequestsRepository_InsertRequest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type input struct {
		userID       string
		sourceID     string
		sourceFormat string
		targetFormat string
		ratio        int
	}

	requestsRepo := NewRequestsRepository(db)
	const query = `insert into converter.requests (.*)`

	testTable := []struct {
		name              string
		args              input
		expectedRequestID string
		mockBehavior      func(input)
		isErrorExpected   bool
	}{
		{
			name: "Ok",
			args: input{
				userID:       "1",
				sourceID:     "1",
				sourceFormat: "jpg",
				targetFormat: "png",
				ratio:        95,
			},
			expectedRequestID: "1",
			mockBehavior: func(args input) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
				mock.ExpectQuery(query).
					WithArgs(args.userID, args.sourceID, args.sourceFormat, args.targetFormat, args.ratio).
					WillReturnRows(rows)
			},
			isErrorExpected: false,
		},
		{
			name: "Database error",
			args: input{
				userID:       "1",
				sourceID:     "1",
				sourceFormat: "jpeg",
				targetFormat: "png",
				ratio:        80,
			},
			mockBehavior: func(args input) {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery(query).
					WithArgs(args.userID, args.sourceID, args.sourceFormat, args.targetFormat, args.ratio).
					WillReturnRows(rows)
			},
			isErrorExpected: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.args)
			result, err := requestsRepo.InsertRequest(context.TODO(),
				tc.args.userID, tc.args.sourceID, tc.args.sourceFormat, tc.args.targetFormat, tc.args.ratio)
			if tc.isErrorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRequestID, result)
			}
		})
	}
}

func TestRequestsRepository_UpdateRequest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type input struct {
		requestID string
		status    string
		targetID  string
	}

	requestsRepo := NewRequestsRepository(db)
	const query = `update converter.requests set (.+)`

	testTable := []struct {
		name            string
		args            input
		expectedError   error
		mockBehavior    func(args input)
		isErrorExpected bool
	}{
		{
			name: "Ok",
			args: input{
				requestID: "1",
				status:    "done",
				targetID:  "123",
			},
			expectedError: nil,
			mockBehavior: func(args input) {
				mock.ExpectExec(query).
					WithArgs(args.requestID, args.targetID, args.status).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			isErrorExpected: false,
		},
		{
			name: "Cannot update request",
			args: input{
				requestID: "1",
				status:    "done",
				targetID:  "123",
			},
			expectedError: fmt.Errorf("can't update request"),
			mockBehavior: func(args input) {
				mock.ExpectExec(query).
					WithArgs(args.requestID, args.targetID, args.status).
					WillReturnError(fmt.Errorf("can't update request"))
			},
			isErrorExpected: true,
		},
		{
			name: "Cannot get affected rows",
			args: input{
				requestID: "1",
				status:    "done",
				targetID:  "123",
			},
			expectedError: fmt.Errorf("can't get the number of rows affected by an update"),
			mockBehavior: func(args input) {
				mock.ExpectExec(query).
					WithArgs(args.requestID, args.targetID, args.status).
					WillReturnError(fmt.Errorf("can't get the number of rows affected by an update"))
			},
			isErrorExpected: true,
		},
		{
			name: "No such request",
			args: input{
				requestID: "1",
				status:    "done",
				targetID:  "123",
			},
			expectedError: ErrNoSuchRequest,
			mockBehavior: func(args input) {
				mock.ExpectExec(query).
					WithArgs(args.requestID, args.targetID, args.status).
					WillReturnError(ErrNoSuchRequest)
			},
			isErrorExpected: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.args)
			err := requestsRepo.UpdateRequest(context.TODO(), tc.args.requestID, tc.args.status, tc.args.targetID)
			if tc.isErrorExpected {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tc.expectedError, err)
				assert.NoError(t, err)
			}
		})
	}
}
