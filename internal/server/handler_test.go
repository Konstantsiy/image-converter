package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	mockservice "github.com/Konstantsiy/image-converter/internal/service/mock"
)

func TestServer_LogIn(t *testing.T) {
	type request struct {
		email    string
		password string
	}

	testTable := []struct {
		name                 string
		requestBody          string
		request              request
		mockBehavior         func(s *mockservice.MockAuthorization, req request)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			requestBody: `
				"email":"email1@gmail.com",
				"password":"password1"`,
			request: request{
				email:    "email1@gmail.com",
				password: "password1",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, req request) {
				s.EXPECT().LogIn(context.TODO(), req.email, req.password).Return(
					LoginResponse{
						AccessToken:  "token1",
						RefreshToken: "token2",
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{
				"access_token"":"token1",
				"refresh_token":"token2"
				}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockservice.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.request)

			//authService := service.NewAuthService(auth, nil)

		})
	}
}
