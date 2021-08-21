// Package hash implements functionality for working with hashed passwords.
package hash

import (
	"golang.org/x/crypto/bcrypt"
)

//GeneratePasswordHash hashes the given password.
func GeneratePasswordHash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash), err
}

//ComparePasswords compares a hashed password with a regular one.
func ComparePasswords(hashedPwd, pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))
	return err != nil, err
}
