package server

import "testing"

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("No error expected but got: %v", err)
	}
}

func assertError(t *testing.T, givenError, expectedError error) {
	t.Helper()
	if givenError != expectedError {
		t.Errorf("Eepected %v, but got another one: %v", expectedError, givenError)
	}
}

// TestValidateUserCredentials tests ValidateUserCredentials function.
func TestValidateUserCredentials(t *testing.T) {
	testTable := []struct {
		Email           string
		Password        string
		IsErrorExpected bool
		ExpectedError   error
	}{
		{
			Email:           "Ivan.Ivanov@company.com",
			Password:        "Simple_Password123",
			IsErrorExpected: false,
			ExpectedError:   nil,
		},
		{
			Email:           "@company.com",
			Password:        "Simple_Password123",
			IsErrorExpected: true,
			ExpectedError:   errInvalidEmailFormat,
		},
		{
			Email:           "@com1111111111",
			Password:        "Simple_Password123",
			IsErrorExpected: true,
			ExpectedError:   errInvalidEmailFormat,
		},
		{
			Email:           "",
			Password:        "Simple_Password123",
			IsErrorExpected: true,
			ExpectedError:   errInvalidEmailLength,
		},
		{
			Email:           "Ivan.Ivanov@company.com",
			Password:        "ae3",
			IsErrorExpected: true,
			ExpectedError:   errInvalidPasswordLength,
		},
		{
			Email:           "Ivan.Ivanov@company.com",
			Password:        "simple_password123",
			IsErrorExpected: true,
			ExpectedError:   errNoOneUppercaseInPassword,
		},
		{
			Email:           "Ivan.Ivanov@company.com",
			Password:        "SIMPLE_PASSWORD123",
			IsErrorExpected: true,
			ExpectedError:   errNoOneLowercaseInPassword,
		},
	}

	for _, tc := range testTable {
		err := ValidateUserCredentials(tc.Email, tc.Password)
		if tc.IsErrorExpected {
			assertError(t, err, tc.ExpectedError)
		} else {
			assertNoError(t, err)
		}
	}
}

// TestValidateConversionRequest tests ValidateConversionRequest function.
func TestValidateConversionRequest(t *testing.T) {
	testTable := []struct {
		Filename        string
		SourceFormat    string
		TargetFormat    string
		Ratio           int
		IsErrorExpected bool
		ExpectedError   error
	}{
		{
			Filename:        "filename",
			SourceFormat:    "jpg",
			TargetFormat:    "png",
			Ratio:           3,
			IsErrorExpected: false,
			ExpectedError:   nil,
		},
		{
			Filename:        "file~name",
			SourceFormat:    "jpg",
			TargetFormat:    "png",
			Ratio:           3,
			IsErrorExpected: true,
			ExpectedError:   errInvalidFilenameFormat,
		},
		{
			Filename:        "filename",
			SourceFormat:    "jpdfdfdg",
			TargetFormat:    "png",
			Ratio:           3,
			IsErrorExpected: true,
			ExpectedError:   errInvalidSourceFormat,
		},
		{
			Filename:        "filename",
			SourceFormat:    "jpg",
			TargetFormat:    "pndssdg",
			Ratio:           3,
			IsErrorExpected: true,
			ExpectedError:   errInvalidTargetFormat,
		},
		{
			Filename:        "filename",
			SourceFormat:    "jpg",
			TargetFormat:    "jpeg",
			Ratio:           3,
			IsErrorExpected: true,
			ExpectedError:   errEqualsFormats,
		},
		{
			Filename:        "filename",
			SourceFormat:    "jpg",
			TargetFormat:    "png",
			Ratio:           -3,
			IsErrorExpected: true,
			ExpectedError:   errInvalidRatio,
		},
	}

	for _, tc := range testTable {
		err := ValidateConversionRequest(tc.Filename, tc.SourceFormat, tc.TargetFormat, tc.Ratio)
		if tc.IsErrorExpected {
			assertError(t, err, tc.ExpectedError)
		} else {
			assertNoError(t, err)
		}
	}
}
