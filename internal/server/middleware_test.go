package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Konstantsiy/image-converter/internal/appcontext"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	mockservice "github.com/Konstantsiy/image-converter/internal/service/mock"
)

func tempHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := appcontext.UserIDFromContext(r.Context())
	if !ok {
		_, _ = w.Write([]byte("cannot get user id"))
	}
	_, _ = w.Write([]byte(userID))
}

func TestServer_AuthMiddleware(t *testing.T) {
	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         func(s *mockservice.MockAuthorization, token string)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mockservice.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return("1", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "1",
		},
		{
			name:                 "Empty header name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(s *mockservice.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: "empty auth handler\n",
		},
		{
			name:                 "Invalid header value",
			headerName:           "Authorization",
			headerValue:          "Ber token",
			token:                "token",
			mockBehavior:         func(s *mockservice.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: "invalid auth handler\n",
		},
		{
			name:                 "Empty token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(s *mockservice.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: "token is empty\n",
		},
		{
			name:        "Token parsing error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mockservice.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return("", fmt.Errorf("invalid token"))
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: "can't parse JWT: invalid token\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockservice.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.token)

			s := &Server{authService: auth}

			r := mux.NewRouter()
			r.Use(s.AuthMiddleware)

			r.HandleFunc("/identity", tempHandler).Methods("GET")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/identity", nil)
			req.Header.Set(tc.headerName, tc.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
