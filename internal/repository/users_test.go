package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewUsersRepository(t *testing.T) {
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
			_, resultErr := NewUsersRepository(tc.mockDB)
			if tc.isErrorExpected {
				assert.Equal(t, resultErr, tc.expectedError)
			} else {
				assert.NoError(t, resultErr)
			}
		})
	}
}

func TestUsersRepository_InsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type input struct {
		email    string
		password string
	}

	usersRepo, err := NewUsersRepository(db)
	require.NoError(t, err)

	const (
		query = "INSERT INTO converter.users (.+) RETURNING id"

		defaultUserUD   = "1"
		defaultEmail    = "email1@gmail.com"
		defaultPassword = "Password1@gmail.com"

		rowUserID = "user_id"
	)

	testTable := []struct {
		name            string
		args            input
		expectedUserID  string
		mockBehavior    func(input)
		isErrorExpected bool
		expectedError   error
	}{
		{
			name: "Ok",
			args: input{
				email:    defaultEmail,
				password: defaultPassword,
			},
			expectedUserID: defaultUserUD,
			mockBehavior: func(args input) {
				rows := sqlmock.NewRows([]string{rowUserID}).AddRow(defaultUserUD)
				mock.ExpectQuery(query).
					WithArgs(args.email, args.password).
					WillReturnRows(rows)
			},
			isErrorExpected: false,
		},
		{
			name: "User already exists",
			args: input{
				email:    defaultEmail,
				password: defaultPassword,
			},
			mockBehavior: func(args input) {
				mock.ExpectQuery(query).
					WithArgs(args.email, args.password).
					WillReturnError(ErrUserAlreadyExists)
			},
			isErrorExpected: true,
			expectedError:   ErrUserAlreadyExists,
		},
		{
			name: "Database error",
			args: input{
				email:    defaultEmail,
				password: defaultPassword,
			},
			mockBehavior: func(args input) {
				rows := sqlmock.NewRows([]string{rowUserID})
				mock.ExpectQuery(query).
					WithArgs(args.email, args.password).
					WillReturnRows(rows)
			},
			isErrorExpected: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.args)
			result, err := usersRepo.InsertUser(context.TODO(), tc.args.email, tc.args.password)
			if tc.isErrorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUserID, result)
			}
		})
	}
}

func TestUsersRepository_GetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	usersRepo, err := NewUsersRepository(db)
	require.NoError(t, err)

	const (
		query = "SELECT (.+) FROM converter.users WHERE (.+)"

		defaultUserID   = "1"
		defaultEmail    = "email1@gmail.com"
		defaultPassword = "Password1"

		rowID       = "ID"
		rowEmail    = "Email"
		rowPassword = "Password"
	)

	testTable := []struct {
		name            string
		email           string
		expectedUser    User
		mockBehavior    func(string)
		isErrorExpected bool
		expectedError   error
	}{
		{
			name:  "Ok",
			email: defaultEmail,
			expectedUser: User{
				ID:       defaultUserID,
				Email:    defaultEmail,
				Password: defaultPassword,
			},
			mockBehavior: func(email string) {
				rows := sqlmock.NewRows([]string{rowID, rowEmail, rowPassword}).
					AddRow("1", "email1@gmail.com", "Password1")
				mock.ExpectQuery(query).
					WithArgs(email).
					WillReturnRows(rows)
			},
			isErrorExpected: false,
		},
		{
			name:  "No such user",
			email: defaultEmail,
			mockBehavior: func(email string) {
				mock.ExpectQuery(query).
					WithArgs(email).
					WillReturnError(ErrNoSuchUser)
			},
			isErrorExpected: true,
			expectedError:   ErrNoSuchUser,
		},
		{
			name:  "Database error",
			email: defaultEmail,
			mockBehavior: func(email string) {
				rows := sqlmock.NewRows([]string{rowID, rowEmail, rowPassword})
				mock.ExpectQuery(query).
					WithArgs(email).
					WillReturnRows(rows)
			},
			isErrorExpected: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.email)
			result, err := usersRepo.GetUserByEmail(context.TODO(), tc.email)
			if tc.isErrorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, result)
			}
		})
	}
}
