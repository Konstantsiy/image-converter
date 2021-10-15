package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Konstantsiy/image-converter/internal/service"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"

	mockservice "github.com/Konstantsiy/image-converter/internal/service/mock"

	"github.com/golang/mock/gomock"
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
			name:        "Ok",
			requestBody: `{"email": "email1@gmail.com","password": "password1"}`,
			request: request{
				email:    "email1@gmail.com",
				password: "password1",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, req request) {
				s.EXPECT().
					LogIn(gomock.Any(), req.email, req.password).
					Return("token1", "token2", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"access_token":"token1","refresh_token":"token2"}`,
		},
		{
			name:        "Invalid input",
			requestBody: `{"email": "","password": "password112"}`,
			request: request{
				email:    "email112@gmail.com",
				password: "password112",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, req request) {
				s.EXPECT().
					LogIn(gomock.Any(), req.email, req.password).
					Return("", "", service.ErrInvalidParam)
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: "invalid email or password\n",
		},
		//{
		//	name:        "Can't generate token access token",
		//	requestBody: `{"email": "email1@gmail.com","password": "password1"}`,
		//	request: request{
		//		email:    "email1@gmail.com",
		//		password: "password1",
		//	},
		//	mockBehavior: func(s *mockservice.MockAuthorization, req request) {
		//		s.EXPECT().
		//			LogIn(gomock.Any(), req.email, req.password).
		//			Return("token1", "token2")
		//	},
		//	expectedStatusCode:   http.StatusInternalServerError,
		//	expectedResponseBody: `"":"token1","refresh_token":"token2"}`,
		//},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockservice.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.request)

			s := Server{auth, nil, nil, nil}

			r := mux.NewRouter()
			r.HandleFunc("/user/login", s.LogIn).Methods("POST")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/user/login", bytes.NewBufferString(tc.requestBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestServer_SignUp(t *testing.T) {
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
			name:        "Ok",
			requestBody: `{"email": "email1@gmail.com","password": "Password1"}`,
			request: request{
				email:    "email1@gmail.com",
				password: "Password1",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, req request) {
				s.EXPECT().
					SignUp(gomock.Any(), req.email, req.password).
					Return("1", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"user_id":"1"}`,
		},
		{
			name:        "Invalid email length",
			requestBody: `{"email": "1@il.cm","password": "Password1"}`,
			request: request{
				email:    "1@il.cm",
				password: "Password1",
			},
			mockBehavior:         func(s *mockservice.MockAuthorization, req request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid email: need minimum 8 characters\n",
		},
		{
			name:        "Invalid email format",
			requestBody: `{"email": "@gmail.com","password": "password1"}`,
			request: request{
				email:    "email1@gmail.com",
				password: "password1",
			},
			mockBehavior:         func(s *mockservice.MockAuthorization, req request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid email: doesn't match the correct address format (for example ivan.ivanov@gmail.com)\n",
		},
		{
			name:        "Invalid password length",
			requestBody: `{"email": "email1@gmail.com","password": "rd1"}`,
			request: request{
				email:    "email1@gmail.com",
				password: "rd1",
			},
			mockBehavior:         func(s *mockservice.MockAuthorization, req request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid password: length must be from 8 to 20 characters\n",
		},
		{
			name:        "Invalid password no lowercase",
			requestBody: `{"email": "email1@gmail.com","password": "PASSWORD1"}`,
			request: request{
				email:    "email1@gmail.com",
				password: "PASSWORD1",
			},
			mockBehavior:         func(s *mockservice.MockAuthorization, req request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid password: need at least one lowercase character\n",
		},
		{
			name:        "Invalid password no uppercase",
			requestBody: `{"email": "email1@gmail.com","password": "password1"}`,
			request: request{
				email:    "email1@gmail.com",
				password: "password1",
			},
			mockBehavior:         func(s *mockservice.MockAuthorization, req request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid password: need at least one uppercase character\n",
		},
		{
			name:        "Invalid password no digit",
			requestBody: `{"email": "email1@gmail.com","password": "Password"}`,
			request: request{
				email:    "email1@gmail.com",
				password: "Password",
			},
			mockBehavior:         func(s *mockservice.MockAuthorization, req request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid password: need at least one digit\n",
		},
		{
			name:        "User already exists",
			requestBody: `{"email": "email1@gmail.com","password": "Password1"}`,
			request: request{
				email:    "email1@gmail.com",
				password: "Password1",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, req request) {
				s.EXPECT().
					SignUp(gomock.Any(), req.email, req.password).
					Return("", &service.ServiceError{
						Err:        errors.New("the user with the given email already exists"),
						StatusCode: http.StatusBadRequest,
					})
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "the user with the given email already exists\n",
		},
		{
			name:        "Cannot generate password hash",
			requestBody: `{"email": "email1@gmail.com","password": "Password1"}`,
			request: request{
				email:    "email1@gmail.com",
				password: "Password1",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, req request) {
				s.EXPECT().
					SignUp(gomock.Any(), req.email, req.password).
					Return("", &service.ServiceError{
						Err:        errors.New("cannot generate password hash"),
						StatusCode: http.StatusInternalServerError,
					})
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: "cannot generate password hash\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockservice.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.request)

			s := Server{auth, nil, nil, nil}

			r := mux.NewRouter()
			r.HandleFunc("/user/signup", s.SignUp).Methods("POST")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/user/signup", bytes.NewBufferString(tc.requestBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
