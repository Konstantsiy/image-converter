// Package hash implements functionality for working with hashed passwords.
package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// GeneratePasswordHash hashes the given password.
func GeneratePasswordHash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash), err
}

// ComparePasswordHash compares a hashed password with a regular one.
func ComparePasswordHash(pwd, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil, err
}
