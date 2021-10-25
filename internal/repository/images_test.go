package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewImagesRepository(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}

	testTable := []struct {
		name            string
		mockDB          *sql.DB
		isErrorExpected bool
		expectedError   error
	}{
		{
			name:            "Ok",
			mockDB:          mockDB,
			isErrorExpected: false,
		},
		{
			name:            "Empty SQL driver",
			mockDB:          nil,
			isErrorExpected: true,
			expectedError:   ErrEmptySQLDriver,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			_, resultErr := NewImagesRepository(tc.mockDB)
			if tc.isErrorExpected {
				assert.Equal(t, resultErr, tc.expectedError)
			} else {
				assert.NoError(t, resultErr)
			}
		})
	}
}

func TestImagesRepository_InsertImage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type input struct {
		filename string
		format   string
	}

	imagesRepo, err := NewImagesRepository(db)
	require.NoError(t, err)

	const query = "INSERT INTO converter.images (.+) RETURNING id"

	testTable := []struct {
		name            string
		args            input
		expectedImageID string
		mockBehavior    func(input)
		isErrorExpected bool
		expectedError   error
	}{
		{
			name: "Ok",
			args: input{
				filename: "image1",
				format:   "jpg",
			},
			expectedImageID: "1",
			mockBehavior: func(args input) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
				mock.ExpectQuery(query).
					WithArgs(args.filename, args.format).
					WillReturnRows(rows)
			},
			isErrorExpected: false,
		},
		{
			name: "Database error",
			args: input{
				filename: "image2",
				format:   "png",
			},
			mockBehavior: func(args input) {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery(query).
					WithArgs(args.filename, args.format).
					WillReturnRows(rows)
			},
			isErrorExpected: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.args)

			result, err := imagesRepo.InsertImage(context.TODO(), tc.args.filename, tc.args.format)

			if tc.isErrorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedImageID, result)
			}
		})
	}
}

func TestImagesRepository_GetImageIDByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type input struct {
		userID  string
		imageID string
	}

	imageRepo, err := NewImagesRepository(db)
	require.NoError(t, err)

	const query = `SELECT (.*)`

	testTable := []struct {
		name            string
		args            input
		expectedImageID string
		mockBehavior    func(input)
		isErrorExpected bool
		expectedError   error
	}{
		{
			name: "Ok",
			args: input{
				userID:  "1",
				imageID: "1",
			},
			expectedImageID: "1",
			mockBehavior: func(args input) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
				mock.ExpectQuery(query).
					WithArgs(args.imageID, args.userID).
					WillReturnRows(rows)
			},
			isErrorExpected: false,
		},
		{
			name: "No such image",
			args: input{
				userID:  "1",
				imageID: "123",
			},
			mockBehavior: func(args input) {
				mock.ExpectQuery(query).
					WithArgs(args.imageID, args.userID).
					WillReturnError(ErrNoSuchImage)
			},
			isErrorExpected: true,
			expectedError:   ErrNoSuchImage,
		},
		{
			name: "Database error",
			args: input{
				userID:  "1",
				imageID: "1",
			},
			mockBehavior: func(args input) {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery(query).
					WithArgs(args.imageID, args.userID).
					WillReturnRows(rows)
			},
			isErrorExpected: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.args)
			result, err := imageRepo.GetImageIDByUserID(context.TODO(), tc.args.userID, tc.args.imageID)
			if tc.isErrorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedImageID, result)
			}
		})
	}
}
