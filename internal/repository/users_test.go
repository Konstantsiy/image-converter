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
	const query = "INSERT INTO converter.users (.+) RETURNING id"

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
				email:    "email1@gmail.com",
				password: "password1@gmail.com",
			},
			expectedUserID: "1",
			mockBehavior: func(args input) {
				rows := sqlmock.NewRows([]string{"user_id"}).AddRow("1")
				mock.ExpectQuery(query).
					WithArgs(args.email, args.password).
					WillReturnRows(rows)
			},
			isErrorExpected: false,
		},
		{
			name: "User already exists",
			args: input{
				email:    "email1@gmail.com",
				password: "password1@gmail.com",
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
				email:    "email1@gmail.com",
				password: "password1@gmail.com",
			},
			mockBehavior: func(args input) {
				rows := sqlmock.NewRows([]string{"id"})
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
	const query = "SELECT (.+) FROM converter.users WHERE (.+)"

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
			email: "email@gmail.com",
			expectedUser: User{
				ID:       "1",
				Email:    "email1@gmail.com",
				Password: "password1",
			},
			mockBehavior: func(email string) {
				rows := sqlmock.NewRows([]string{"ID", "Email", "Password"}).
					AddRow("1", "email1@gmail.com", "password1")
				mock.ExpectQuery(query).
					WithArgs(email).
					WillReturnRows(rows)
			},
			isErrorExpected: false,
		},
		{
			name:  "No such user",
			email: "email@gmail.com",
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
			email: "email@gmail.com",
			mockBehavior: func(email string) {
				rows := sqlmock.NewRows([]string{"ID", "Email", "Password"})
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
