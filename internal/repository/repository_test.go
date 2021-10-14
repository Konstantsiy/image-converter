package repository

//
//import (
//	"database/sql"
//	"errors"
//	"testing"
//
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/gofrs/uuid"
//	//"github.com/gofrs/uuid"
//	"github.com/lib/pq"
//)
//
////// NewMock creates new sqlmock instance.
////func NewMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
////	db, mock, err := sqlmock.New()
////	if err != nil {
////		t.Errorf("an error '%v' was not expected when opening a stub database connection", err)
////	}
////
////	return db, mock
////}
//
//func TestRepository_InsertUser(t *testing.T) {
//	db, mock := NewMock(t)
//	defer db.Close()
//
//	type input struct {
//		email    string
//		password string
//	}
//
//	repo := NewRepository(db)
//	const query = "insert into converter.users (.*) returning id"
//
//	testTable := []struct {
//		name           string
//		input          input
//		expectedUserID string
//		expectedError  error
//		mockBehavior   func(input, string)
//	}{
//		{
//			name: "Ok",
//			input: input{
//				email:    "email2",
//				password: "password2",
//			},
//			expectedUserID: uuid.NewV4().String(),
//			expectedError:  nil,
//			mockBehavior: func(input input, userID string) {
//				mock.ExpectQuery(query).
//					WithArgs(input.email, input.password).
//					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))
//			},
//		},
//		{
//			name:          "User already exists",
//			expectedError: ErrUserAlreadyExists,
//			mockBehavior: func(input input, userID string) {
//				mock.ExpectQuery(query).
//					WithArgs(input.email, input.password).
//					WillReturnError(&pq.Error{Code: uniqueViolationCode})
//			},
//		},
//	}
//
//	for _, tc := range testTable {
//		t.Run(tc.name, func(t *testing.T) {
//			tc.mockBehavior(tc.input, tc.expectedUserID)
//
//			resultUserID, err := repo.LogIn(tc.input.email, tc.input.password)
//			if !errors.Is(err, tc.expectedError) {
//				t.Errorf("expected error to be %v, but got %v", tc.expectedError, err)
//			}
//
//			if tc.expectedUserID != resultUserID {
//				t.Errorf("expected user id to be %q, but got %q", tc.expectedUserID, resultUserID)
//			}
//
//			if err := mock.ExpectationsWereMet(); err != nil {
//				t.Errorf("there were unfulfilled expectations: %v", err)
//			}
//		})
//	}
//}
//
//func TestRepository_GetUserByEmail(t *testing.T) {
//	db, mock := NewMock(t)
//	defer db.Close()
//
//	repo := NewRepository(db)
//	const query = "select (.*) from converter.users where email = ?"
//
//	testTable := []struct {
//		name          string
//		email         string
//		expectedUser  User
//		expectedError error
//		mockBehavior  func(string, User)
//	}{
//		{
//			name:  "Ok",
//			email: "email1",
//			expectedUser: User{
//				ID:       "1",
//				Email:    "email1",
//				Password: "password1",
//			},
//			expectedError: nil,
//			mockBehavior: func(email string, user User) {
//				mock.ExpectQuery(query).
//					WithArgs(email).
//					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
//						AddRow(user.ID, user.Email, user.Password))
//			},
//		},
//		{
//			name:          "No such user",
//			email:         "email23234",
//			expectedError: ErrNoSuchUser,
//			mockBehavior: func(email string, user User) {
//				mock.ExpectQuery(query).
//					WithArgs(email).
//					WillReturnError(sql.ErrNoRows)
//			},
//		},
//	}
//
//	for _, tc := range testTable {
//		t.Run(tc.name, func(t *testing.T) {
//			tc.mockBehavior(tc.email, tc.expectedUser)
//
//			resultUser, err := repo.SignUp(tc.email)
//			if !errors.Is(err, tc.expectedError) {
//				t.Errorf("expected error to be %v, but got %v", tc.expectedError, err)
//			}
//
//			if tc.expectedUser != resultUser {
//				t.Errorf("expected user id to be %v, but got %v", tc.expectedUser, resultUser)
//			}
//
//			if err := mock.ExpectationsWereMet(); err != nil {
//				t.Errorf("there were unfulfilled expectations: %v", err)
//			}
//		})
//	}
//}
