// Package validation implements functions that validate requests data.
package validation

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	// emailRegex checks the validity of the format of the entire email address.
	emailRegex = "^[A-Za-z0-9._]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$"

	// passwordOneLowercaseRegex searches for at least one lowercase character in password.
	passwordOneLowercaseRegex = "(.*[a-z])"

	// passwordOneUppercaseRegex searches for at least one uppercase character in password.
	passwordOneUppercaseRegex = "(.*[A-Z])"

	// passwordOneDigitRegex searches for at least one digit in password.
	passwordOneDigitRegex = "(.*[0-9])"

	// filenameInvalidCharactersRegex searches for special characters in password.
	filenameInvalidCharactersRegex = "[:#%^;<>\\{}[\\]+~`=?&,\" ]"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 20
	minEmailLength    = 8
	minRatio          = 1
	maxRatio          = 9
)

var formats = map[string]struct{}{
	"jpg": {},
	"png": {},
}

var (
	errInvalidEmailLength = fmt.Errorf("email has minimum %d characters", minPasswordLength)
	errInvalidEmailFormat = errors.New("invalid email address format")

	errInvalidPasswordLength    = fmt.Errorf("password length must be from %d to %d characters", minPasswordLength, maxPasswordLength)
	errNoOneLowercaseInPassword = errors.New("password must contain at least one lowercase character")
	errNoOneUppercaseInPassword = errors.New("password must contain at least one uppercase character")
	errNoOneDigitInPassword     = errors.New("password must contain at least one digit")

	errMissingFilename       = errors.New("missing filename")
	errInvalidFilenameFormat = errors.New("filename shouldn't contain space and any special characters like :;<>{}[]+=?&,\"")

	errInvalidSourceFormat = errors.New("invalid source format: needed jpg or png")
	errInvalidTargetFormat = errors.New("invalid target format: needed jpg or png")
	errEqualsFormats       = errors.New("source and target formats should differ")

	errInvalidRatio = fmt.Errorf("invalid ratio: needed a value from %d to %d inclusive", minRatio, maxRatio)
)

// ValidateSignUpRequest validates user credentials.
func ValidateSignUpRequest(email, password string) error {
	if len(email) < minEmailLength {
		return errInvalidEmailLength
	}

	if match, _ := regexp.MatchString(emailRegex, email); !match {
		return errInvalidEmailFormat
	}

	if l := len(password); l < minPasswordLength || l > maxPasswordLength {
		return errInvalidPasswordLength
	}

	if match, _ := regexp.MatchString(passwordOneLowercaseRegex, password); !match {
		return errNoOneLowercaseInPassword
	}

	if match, _ := regexp.MatchString(passwordOneUppercaseRegex, password); !match {
		return errNoOneUppercaseInPassword
	}

	if match, _ := regexp.MatchString(passwordOneDigitRegex, password); !match {
		return errNoOneDigitInPassword
	}

	return nil
}

// ValidateConversionRequest validates data from the conversion request body.
func ValidateConversionRequest(filename, sourceFormat, targetFormat string, ratio int) error {
	if len(filename) == 0 {
		return errMissingFilename
	}

	if match, _ := regexp.MatchString(filenameInvalidCharactersRegex, filename); match {
		return errInvalidFilenameFormat
	}

	if _, ok := formats[sourceFormat]; !ok {
		return errInvalidSourceFormat
	}

	if _, ok := formats[targetFormat]; !ok {
		return errInvalidTargetFormat
	}

	if sourceFormat == targetFormat {
		return errEqualsFormats
	}

	if ratio < minRatio || ratio > maxRatio {
		return errInvalidRatio
	}

	return nil
}
