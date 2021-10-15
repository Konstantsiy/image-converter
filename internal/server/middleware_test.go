package server

//
//import (
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/gorilla/mux"
//	"github.com/stretchr/testify/assert"
//
//	"github.com/golang/mock/gomock"
//
//	mockservice "github.com/Konstantsiy/image-converter/internal/service/mock"
//)
//
//func TestServer_AuthMiddleware(t *testing.T) {
//	testTable := []struct {
//		name                 string
//		headerName           string
//		headerValue          string
//		token                string
//		mockBehavior         func(s *mockservice.MockAuthorization, token string)
//		expectedStatusCode   int
//		expectedResponseBody string
//	}{
//		{
//			name:        "Ok",
//			headerName:  "Authorization",
//			headerValue: "Bearer token",
//			token:       "token",
//			mockBehavior: func(s *mockservice.MockAuthorization, token string) {
//				s.EXPECT().ParseToken(token).Return("1", nil)
//			},
//			expectedStatusCode:   http.StatusOK,
//			expectedResponseBody: "1",
//		},
//	}
//
//	for _, tc := range testTable {
//		t.Run(tc.name, func(t *testing.T) {
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			auth := mockservice.NewMockAuthorization(c)
//			tc.mockBehavior(auth, tc.token)
//
//			s := Server{authService: auth}
//
//			r := mux.NewRouter()
//			r.HandleFunc("/authorization", s.LogIn)
//
//			w := httptest.NewRecorder()
//			req := httptest.NewRequest("GET", "/authorization", nil)
//			req.Header.Set(tc.headerName, tc.headerValue)
//
//			r.ServeHTTP(w, req)
//
//			assert.Equal(t, tc.expectedStatusCode, w.Code)
//			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
//		})
//	}
//}
