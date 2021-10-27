package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	// ErrNoSuchUser notifies that the needed user does not exist.
	ErrNoSuchUser = errors.New("the user with this email does not exist")

	// ErrUserAlreadyExists notifies that a user with such an email already exists.
	ErrUserAlreadyExists = errors.New("the user with the given email already exists")
)

// uniqueViolationCode represents an error code that such an entity already exists.
const uniqueViolationCode = "23505"

// User represents the user in the database.
type User struct {
	ID       string
	Email    string
	Password string
}

// UsersRepository represents repository fro working with users.
type UsersRepository struct {
	db *sql.DB
}

// NewUsersRepository creates new users repository.
func NewUsersRepository(db *sql.DB) (*UsersRepository, error) {
	if db == nil {
		return nil, ErrEmptySQLDriver
	}
	return &UsersRepository{db: db}, nil
}

// InsertUser inserts the user into users table and returns user id.
func (ur *UsersRepository) InsertUser(ctx context.Context, email, password string) (string, error) {
	var userID string
	const query = "INSERT INTO converter.users (email, password) VALUES ($1, $2) RETURNING id;"

	err := ur.db.QueryRowContext(ctx, query, email, password).Scan(&userID)
	if err, ok := err.(*pq.Error); ok && err.Code == uniqueViolationCode {
		return "", ErrUserAlreadyExists
	}
	if err != nil {
		return "", fmt.Errorf("can't insert user: %w", err)
	}

	return userID, nil
}

// GetUserByEmail gets the information about the user by given email.
func (ur *UsersRepository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	const query = "SELECT id, email, password FROM converter.users WHERE email = $1;"

	err := ur.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return User{}, ErrNoSuchUser
	}
	if err != nil {
		return User{}, fmt.Errorf("error in the user selection: %w", err)
	}

	return user, nil
}
