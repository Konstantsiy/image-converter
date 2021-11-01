package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Konstantsiy/image-converter/internal/repository"

	"github.com/Konstantsiy/image-converter/internal/service"
	mockservice "github.com/Konstantsiy/image-converter/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
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
			requestBody: `{"email": "email1@gmail.com","password": "Password1"}`,
			request: request{
				email:    "email1@gmail.com",
				password: "Password1",
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
			name:        "Empty email",
			requestBody: `{"email": "","password": "Password1"}`,
			request: request{
				email:    "",
				password: "Password1",
			},
			mockBehavior:         func(s *mockservice.MockAuthorization, req request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "empty field\n",
		},
		{
			name:        "Empty password",
			requestBody: `{"email": "email1@gmail.com","password": ""}`,
			request: request{
				email:    "email1@gmail.com",
				password: "",
			},
			mockBehavior:         func(s *mockservice.MockAuthorization, req request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "empty field\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockservice.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.request)

			s := Server{authService: auth}

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
						Err:        fmt.Errorf("the user with the given email already exists"),
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
						Err:        fmt.Errorf("cannot generate password hash"),
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

			s := Server{authService: auth}

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

func TestServer_DownloadImage(t *testing.T) {
	const defaultImageID = "1"

	testTable := []struct {
		name                 string
		imageID              string
		mockBehavior         func(s *mockservice.MockImages, id string)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Ok",
			imageID: defaultImageID,
			mockBehavior: func(s *mockservice.MockImages, id string) {
				s.EXPECT().
					Download(gomock.Any(), id).
					Return(defaultImageID, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"image_url":"1"}`,
		},
		{
			name:                 "Missing image id",
			imageID:              "",
			mockBehavior:         func(s *mockservice.MockImages, id string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "image id is missing in parameters\n",
		},
		{
			name:    "Cannot get user id from context",
			imageID: defaultImageID,
			mockBehavior: func(s *mockservice.MockImages, id string) {
				s.EXPECT().
					Download(gomock.Any(), id).
					Return("", &service.ServiceError{
						Err:        fmt.Errorf("can't get user id from application context"),
						StatusCode: http.StatusUnauthorized,
					})
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: "can't get user id from application context\n",
		},
		{
			name:    "No such image",
			imageID: defaultImageID,
			mockBehavior: func(s *mockservice.MockImages, id string) {
				s.EXPECT().
					Download(gomock.Any(), id).
					Return("", &service.ServiceError{
						Err:        fmt.Errorf("no such image"),
						StatusCode: http.StatusNotFound,
					})
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: "no such image\n",
		},
		{
			name:    "Storage error",
			imageID: defaultImageID,
			mockBehavior: func(s *mockservice.MockImages, id string) {
				s.EXPECT().
					Download(gomock.Any(), id).
					Return("", &service.ServiceError{
						Err:        fmt.Errorf("can't get image url"),
						StatusCode: http.StatusInternalServerError,
					})
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: "can't get image url\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			is := mockservice.NewMockImages(c)
			tc.mockBehavior(is, tc.imageID)

			s := Server{imageService: is}

			r := mux.NewRouter()
			r.HandleFunc("/images", s.DownloadImage).Methods("GET")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/images", nil)

			q := req.URL.Query()
			q.Add("id", tc.imageID)
			req.URL.RawQuery = q.Encode()

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestServer_GetRequestsHistory(t *testing.T) {
	defaultTime := time.Now()
	defaultResponseBody := []repository.ConversionRequest{
		{ID: "1", UserID: "1", SourceID: "11", TargetID: "12", SourceFormat: "jpg", TargetFormat: "png",
			Ratio: 90, Created: defaultTime, Updated: defaultTime, Status: "done"},
		{ID: "2", UserID: "2", SourceID: "22", TargetID: "23", SourceFormat: "jpeg", TargetFormat: "png",
			Ratio: 90, Created: defaultTime, Updated: defaultTime, Status: "processing"},
	}

	respJSON, err := json.Marshal(defaultResponseBody)
	if err != nil {
		t.Fatalf("can't marshal response body: %v", err)
	}

	testTable := []struct {
		name                 string
		mockBehavior         func(s *mockservice.MockRequests)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(s *mockservice.MockRequests) {
				s.EXPECT().
					GetUsersRequests(gomock.Any()).
					Return(defaultResponseBody, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: string(respJSON),
		},
		{
			name: "Cannot get user id from context",
			mockBehavior: func(s *mockservice.MockRequests) {
				s.EXPECT().
					GetUsersRequests(gomock.Any()).
					Return(nil, &service.ServiceError{
						Err:        fmt.Errorf("can't get user id from application context"),
						StatusCode: http.StatusInternalServerError,
					})
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: "can't get user id from application context\n",
		},
		{
			name: "Repository error",
			mockBehavior: func(s *mockservice.MockRequests) {
				s.EXPECT().
					GetUsersRequests(gomock.Any()).
					Return(nil, &service.ServiceError{
						Err:        fmt.Errorf("repository error"),
						StatusCode: http.StatusInternalServerError,
					})
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: "repository error\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			rs := mockservice.NewMockRequests(c)
			tc.mockBehavior(rs)

			s := Server{requestsService: rs}

			r := mux.NewRouter()
			r.HandleFunc("/requests", s.GetRequestsHistory).Methods("GET")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/requests", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func createMockRequest(t *testing.T, filename, formFileKey, url, method string, params map[string]string) *http.Request {
	file, err := os.Create(filename)
	require.NoError(t, err)
	defer func() {
		err := file.Close()
		require.NoError(t, err)
	}()

	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	for x := 0; x < 20; x++ {
		for y := 0; y < 20; y++ {
			img.Set(x, y, color.White)
		}
	}

	err = jpeg.Encode(file, img, nil)
	require.NoError(t, err)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(formFileKey, filepath.Base(filename))
	require.NoError(t, err)
	_, _ = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	require.NoError(t, err)

	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)

	req.Header.Set(ContentTypeKey, writer.FormDataContentType())

	return req
}

func TestServer_ConvertImage(t *testing.T) {
	const (
		defaultFormFile     = "file"
		defaultFilename     = "Screenshot_1.jpg"
		targetFormatKey     = "targetFormat"
		defaultTargetFormat = "png"
		defaultRatio        = "90"
		ratioKey            = "ratio"
	)

	type request struct {
		formFileKey string
		filename    string
		params      map[string]string
	}

	testTable := []struct {
		name                 string
		request              request
		mockBehavior         func(s *mockservice.MockImages, p *mockservice.MockProducer, request request)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			request: request{
				formFileKey: defaultFormFile,
				filename:    defaultFilename,
				params: map[string]string{
					targetFormatKey: defaultTargetFormat,
					ratioKey:        defaultRatio,
				},
			},
			mockBehavior: func(s *mockservice.MockImages, p *mockservice.MockProducer, request request) {
				s.EXPECT().
					Convert(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
						gomock.Any(), gomock.Any()).Return("1", "1", nil)
				p.EXPECT().
					SendToQueue(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedStatusCode:   http.StatusAccepted,
			expectedResponseBody: `{"request_id":"1"}`,
		},
		{
			name: "Invalid file form",
			request: request{
				formFileKey: "fill",
				filename:    defaultFilename,
				params: map[string]string{
					targetFormatKey: defaultTargetFormat,
					ratioKey:        defaultRatio,
				},
			},
			mockBehavior:         func(s *mockservice.MockImages, p *mockservice.MockProducer, request request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "can't get sourceFile from form\n",
		},
		{
			name: "Invalid ration form value",
			request: request{
				formFileKey: defaultFormFile,
				filename:    defaultFilename,
				params: map[string]string{
					targetFormatKey: defaultTargetFormat,
					ratioKey:        "ratio_string_value",
				},
			},
			mockBehavior:         func(s *mockservice.MockImages, p *mockservice.MockProducer, request request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid ratio form value\n",
		},
		{
			name: "Invalid filename_1",
			request: request{
				formFileKey: defaultFormFile,
				filename:    "Screenshot_?.jpg",
				params: map[string]string{
					targetFormatKey: defaultTargetFormat,
					ratioKey:        defaultRatio,
				},
			},
			mockBehavior:         func(s *mockservice.MockImages, p *mockservice.MockProducer, request request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid filename: shouldn't contain space and any special characters like :;<>{}[]+=?&,\"\n",
		},
		{
			name: "Invalid filename_2",
			request: request{
				formFileKey: defaultFormFile,
				filename:    "Screensho.t_123.jpg",
				params: map[string]string{
					targetFormatKey: defaultTargetFormat,
					ratioKey:        defaultRatio,
				},
			},
			mockBehavior:         func(s *mockservice.MockImages, p *mockservice.MockProducer, request request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid source format: needed jpg or png\n",
		},
		{
			name: "Invalid source format",
			request: request{
				formFileKey: defaultFormFile,
				filename:    "Screensho.t_1jpg",
				params: map[string]string{
					targetFormatKey: defaultTargetFormat,
					ratioKey:        defaultRatio,
				},
			},
			mockBehavior:         func(s *mockservice.MockImages, p *mockservice.MockProducer, request request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid source format: needed jpg or png\n",
		},
		{
			name: "Invalid target format",
			request: request{
				formFileKey: defaultFormFile,
				filename:    defaultFilename,
				params: map[string]string{
					targetFormatKey: "pngdfdf",
					ratioKey:        defaultRatio,
				},
			},
			mockBehavior:         func(s *mockservice.MockImages, p *mockservice.MockProducer, request request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid target format: needed jpg or png\n",
		},
		{
			name: "Invalid formats",
			request: request{
				formFileKey: defaultFormFile,
				filename:    defaultFilename,
				params: map[string]string{
					targetFormatKey: "jpg",
					ratioKey:        defaultRatio,
				},
			},
			mockBehavior:         func(s *mockservice.MockImages, p *mockservice.MockProducer, request request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid formats: source and target formats should differ\n",
		},
		{
			name: "Invalid ratio",
			request: request{
				formFileKey: defaultFormFile,
				filename:    defaultFilename,
				params: map[string]string{
					targetFormatKey: defaultTargetFormat,
					ratioKey:        "101",
				},
			},
			mockBehavior:         func(s *mockservice.MockImages, p *mockservice.MockProducer, request request) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid ratio: needed a value from 1 to 99 inclusive\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			is := mockservice.NewMockImages(c)
			p := mockservice.NewMockProducer(c)

			s := Server{imageService: is, producer: p}

			r := mux.NewRouter()
			r.HandleFunc("/conversion", s.ConvertImage).Methods("POST")

			w := httptest.NewRecorder()

			req := createMockRequest(t, tc.request.filename, tc.request.formFileKey, "/conversion", "POST",
				tc.request.params)
			defer os.Remove(tc.request.filename)

			tc.mockBehavior(is, p, tc.request)

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
