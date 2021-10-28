// Package validation implements functions that validate requests data.
package validation

import (
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
	filenameInvalidCharactersRegex = "[:#%^;<>\\{}[\\]+~`=.?&,\" ]"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 20
	minEmailLength    = 8
	minRatio          = 1
	maxRatio          = 99
)

var formats = map[string]struct{}{
	"jpg": {},
	"png": {},
}

// InvalidParameterError represents validation related error.
type InvalidParameterError struct {
	Param   string
	Message string
}

func (e *InvalidParameterError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Param, e.Message)
}

// ValidateSignUpRequest validates user credentials.
func ValidateSignUpRequest(email, password string) error {
	if len(email) < minEmailLength {
		return &InvalidParameterError{
			Param:   "email",
			Message: fmt.Sprintf("need minimum %d characters", minEmailLength),
		}
	}

	if match, _ := regexp.MatchString(emailRegex, email); !match {
		return &InvalidParameterError{
			Param:   "email",
			Message: "doesn't match the correct address format (for example ivan.ivanov@gmail.com)",
		}
	}

	if l := len(password); l < minPasswordLength || l > maxPasswordLength {
		return &InvalidParameterError{
			Param:   "password",
			Message: fmt.Sprintf("length must be from %d to %d characters", minPasswordLength, maxPasswordLength),
		}
	}

	if match, _ := regexp.MatchString(passwordOneLowercaseRegex, password); !match {
		return &InvalidParameterError{
			Param:   "password",
			Message: "need at least one lowercase character",
		}
	}

	if match, _ := regexp.MatchString(passwordOneUppercaseRegex, password); !match {
		return &InvalidParameterError{
			Param:   "password",
			Message: "need at least one uppercase character",
		}
	}

	if match, _ := regexp.MatchString(passwordOneDigitRegex, password); !match {
		return &InvalidParameterError{
			Param:   "password",
			Message: "need at least one digit",
		}
	}

	return nil
}

// ValidateConversionRequest validates data from the conversion request body.
func ValidateConversionRequest(filename, sourceFormat, targetFormat string, ratio int) error {
	if filename == "" {
		return &InvalidParameterError{
			Param:   "filename",
			Message: "shouldn't be empty",
		}
	}

	if match, _ := regexp.MatchString(filenameInvalidCharactersRegex, filename); match {
		return &InvalidParameterError{
			Param:   "filename",
			Message: "shouldn't contain space and any special characters like :;<>{}[]+=?&,\"",
		}
	}

	if targetFormat == "jpeg" {
		targetFormat = "jpg"
	}
	if sourceFormat == "jpeg" {
		sourceFormat = "jpg"
	}

	if _, ok := formats[sourceFormat]; !ok {
		return &InvalidParameterError{
			Param:   "source format",
			Message: "needed jpg or png",
		}
	}

	if _, ok := formats[targetFormat]; !ok {
		return &InvalidParameterError{
			Param:   "target format",
			Message: "needed jpg or png",
		}
	}

	if sourceFormat == targetFormat {
		return &InvalidParameterError{
			Param:   "formats",
			Message: "source and target formats should differ",
		}
	}

	if ratio < minRatio || ratio > maxRatio {
		return &InvalidParameterError{
			Param:   "ratio",
			Message: fmt.Sprintf("needed a value from %d to %d inclusive", minRatio, maxRatio),
		}
	}

	return nil
}
