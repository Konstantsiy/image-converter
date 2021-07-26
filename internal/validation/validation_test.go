package validation

import (
	"testing"
)

func assertNoError(t *testing.T, param string) {
	t.Helper()
	if param != "" {
		t.Errorf("No error expected but got invalid %s", param)
	}
}

func assertError(t *testing.T, givenParam, expectedParam string) {
	t.Helper()
	if givenParam != expectedParam {
		t.Errorf("Expected invalid %s, but got invalid %s", givenParam, expectedParam)
	}
}

func TestValidateSignUpRequest(t *testing.T) {
	testTable := []struct {
		Email                string
		Password             string
		IsErrorExpected      bool
		ExpectedInvalidParam string
	}{
		{
			Email:                "Ivan.Ivanov@company.com",
			Password:             "Simple_Password123",
			IsErrorExpected:      false,
			ExpectedInvalidParam: "",
		},
		{
			Email:                "@company.com",
			Password:             "Simple_Password123",
			IsErrorExpected:      true,
			ExpectedInvalidParam: "email",
		},
		{
			Email:                "@com1111111111",
			Password:             "Simple_Password123",
			IsErrorExpected:      true,
			ExpectedInvalidParam: "email",
		},
		{
			Email:                "",
			Password:             "Simple_Password123",
			IsErrorExpected:      true,
			ExpectedInvalidParam: "email",
		},
		{
			Email:                "Ivan.Ivanov@company.com",
			Password:             "ae3",
			IsErrorExpected:      true,
			ExpectedInvalidParam: "password",
		},
		{
			Email:                "Ivan.Ivanov@company.com",
			Password:             "simple_password123",
			IsErrorExpected:      true,
			ExpectedInvalidParam: "password",
		},
		{
			Email:                "Ivan.Ivanov@company.com",
			Password:             "SIMPLE_PASSWORD123",
			IsErrorExpected:      true,
			ExpectedInvalidParam: "password",
		},
	}

	for _, tc := range testTable {
		err := ValidateSignUpRequest(tc.Email, tc.Password)
		if tc.IsErrorExpected {
			assertError(t, err.Param, tc.ExpectedInvalidParam)
		} else {
			assertNoError(t, err.Param)
		}
	}
}

func TestValidateConversionRequest(t *testing.T) {
	testTable := []struct {
		Filename             string
		SourceFormat         string
		TargetFormat         string
		Ratio                int
		IsErrorExpected      bool
		ExpectedInvalidParam string
	}{
		{
			Filename:             "filename",
			SourceFormat:         "jpg",
			TargetFormat:         "png",
			Ratio:                3,
			IsErrorExpected:      false,
			ExpectedInvalidParam: "",
		},
		{
			Filename:             "file~name",
			SourceFormat:         "jpg",
			TargetFormat:         "png",
			Ratio:                3,
			IsErrorExpected:      true,
			ExpectedInvalidParam: "filename",
		},
		{
			Filename:             "filename",
			SourceFormat:         "jpdfdfdg",
			TargetFormat:         "png",
			Ratio:                3,
			IsErrorExpected:      true,
			ExpectedInvalidParam: "source format",
		},
		{
			Filename:             "filename",
			SourceFormat:         "jpg",
			TargetFormat:         "pndssdg",
			Ratio:                3,
			IsErrorExpected:      true,
			ExpectedInvalidParam: "target format",
		},
		{
			Filename:             "filename",
			SourceFormat:         "jpg",
			TargetFormat:         "jpg",
			Ratio:                3,
			IsErrorExpected:      true,
			ExpectedInvalidParam: "formats",
		},
		{
			Filename:             "filename",
			SourceFormat:         "jpg",
			TargetFormat:         "png",
			Ratio:                -3,
			IsErrorExpected:      true,
			ExpectedInvalidParam: "ratio",
		},
	}

	for _, tc := range testTable {
		err := ValidateConversionRequest(tc.Filename, tc.SourceFormat, tc.TargetFormat, tc.Ratio)
		if tc.IsErrorExpected {
			assertError(t, err.Param, tc.ExpectedInvalidParam)
		} else {
			assertNoError(t, err.Param)
		}
	}
}
