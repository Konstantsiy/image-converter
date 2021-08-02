package validation

import (
	"testing"
)

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("No error expected but got: %v", err)
	}
}

func assertError(t *testing.T, givenParam, expectedParam string) {
	t.Helper()
	if givenParam != expectedParam {
		t.Errorf("Expected invalid %s, but got invalid %s", expectedParam, givenParam)
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
		verr, _ := err.(*InvalidParameterError)
		if tc.IsErrorExpected {
			assertError(t, verr.Param, tc.ExpectedInvalidParam)
		} else {
			assertNoError(t, err)
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
		verr, _ := err.(*InvalidParameterError)
		if tc.IsErrorExpected {
			assertError(t, verr.Param, tc.ExpectedInvalidParam)
		} else {
			assertNoError(t, err)
		}
	}
}
