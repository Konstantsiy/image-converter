package server

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
	minLength = 8
	maxLength = 20

	minRatio = 1
	maxRatio = 9
)

var formats = []string{"jpeg", "jpg", "png"}

var (
	errInvalidEmailLength = fmt.Errorf("invalid email length: email it must be at least %d characters", minLength)
	errInvalidEmailFormat = errors.New("invalid email format: enter the correct format, for example Ivan.Ivanov@company.com")

	errInvalidPasswordLength    = fmt.Errorf("invalid password length: password must be from %d to %d characters", minLength, maxLength)
	errNoOneLowercaseInPassword = errors.New("password must contain at least one lowercase character")
	errNoOneUppercaseInPassword = errors.New("the password must contain at least one uppercase character")
	errNoOneDigitInPassword     = errors.New("the password must contain at least one digit")

	errMissingFilename       = errors.New("missing filename")
	errInvalidFilenameFormat = errors.New("filename shouldn't contain space and any special characters like :;<>{}[]+=?&,\"")

	errInvalidSourceFormat = errors.New("invalid source format: needed jpeg, jpg or png")
	errInvalidTargetFormat = errors.New("invalid source format: needed jpeg, jpg or png")
	errEqualsFormats       = errors.New("source and target formats should differ")

	errInvalidRatio = fmt.Errorf("invalid compressoin ratio value: needed a value from %d to %d inclusive", minRatio, maxRatio)
)

// ValidateUserCredentials validates user credentials.
func ValidateUserCredentials(email, password string) error {
	if len(email) < minLength {
		return errInvalidEmailLength
	}

	if match, _ := regexp.MatchString(emailRegex, email); !match {
		return errInvalidEmailFormat
	}

	if l := len(password); l < minLength || l > maxLength {
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

	if !isFormatExists(sourceFormat) {
		return errInvalidSourceFormat
	}

	if !isFormatExists(targetFormat) {
		return errInvalidTargetFormat
	}

	if sourceFormat == targetFormat ||
		(sourceFormat == "jpeg" && targetFormat == "jpg") ||
		(sourceFormat == "jpg" && targetFormat == "jpeg") {
		return errEqualsFormats
	}

	if ratio < minRatio || ratio > maxRatio {
		return errInvalidRatio
	}

	return nil
}

// isFormatExists checks the existence of the needed format.
func isFormatExists(format string) bool {
	for _, item := range formats {
		if item == format {
			return true
		}
	}
	return false
}
