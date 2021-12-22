package tests

import (
	"bytes"
	"context"
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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/pkg/hash"
)

const (
	testLoginURL      = "/user/login"
	testSignUpURL     = "/user/signup"
	testDownloadURL   = "/images"
	testConversionURL = "/conversion"
	testHistoryURL    = "/requests"

	headerTypeKey   = "Content-type"
	headerTypeValue = "application/json"
	headerAuthKey   = "Authorization"
	headerAuthValue = "Bearer "

	defaultEmail        = "email3@gmail.com"
	defaultPassword     = "Password3"
	defaultFile         = "file1"
	defaultFilename     = "Screenshot_1.jpg"
	defaultSourceFormat = "jpg"
	defaultTargetFormat = "png"
	defaultRatio        = "99"
	defaultImageURL     = "123"
)

func (s *APITestSuite) truncateTableUsers() {
	query := "TRUNCATE TABLE converter.users CASCADE;"
	_, err := s.db.Exec(query)
	if err != nil {
		s.FailWithError(fmt.Errorf("unable to truncate users table: %v", err))
	}
}

func (s *APITestSuite) TestUserSignIn() {
	signInData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, defaultEmail, defaultPassword)
	pwdHash, err := hash.GeneratePasswordHash(defaultPassword)
	s.NoError(err)

	_, err = s.repos.users.InsertUser(context.Background(), defaultEmail, pwdHash)
	s.NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, testLoginURL, bytes.NewBuffer([]byte(signInData)))
	req.Header.Set(headerTypeKey, headerTypeValue)

	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.truncateTableUsers()
}

func (s *APITestSuite) TestUserSignUp() {
	signUpData := fmt.Sprintf(`{"email" :"%s","password": "%s"}`, defaultEmail, defaultPassword)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, testSignUpURL, bytes.NewBuffer([]byte(signUpData)))
	req.Header.Set(headerTypeKey, headerTypeValue)

	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusCreated, w.Result().StatusCode)

	type response struct {
		UserID string `json:"user_id"`
	}
	var resp response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)

	ur, err := repository.NewUsersRepository(s.db)
	s.NoError(err)

	user, err := ur.GetUserByEmail(context.Background(), defaultEmail)
	s.NoError(err)

	s.Equal(defaultEmail, user.Email)
	equal, err := hash.ComparePasswordHash(defaultPassword, user.Password)
	s.NoError(err)
	s.True(equal)
	s.truncateTableUsers()
}

func (s *APITestSuite) TestGetRequestsHistory() {
	userID, err := s.repos.users.InsertUser(context.Background(), defaultEmail, defaultPassword)
	s.NoError(err)

	sourceID, err := s.repos.images.InsertImage(context.Background(), defaultFile, defaultSourceFormat)
	s.NoError(err)

	rr, err := repository.NewRequestsRepository(s.db)
	s.NoError(err)
	requestID, err := rr.InsertRequest(context.Background(), userID, sourceID, defaultSourceFormat, defaultTargetFormat, 99)
	s.NoError(err)

	jwt, err := s.tm.GenerateAccessToken(userID)
	s.NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, testHistoryURL, nil)
	req.Header.Set(headerTypeKey, headerTypeValue)
	req.Header.Set(headerAuthKey, headerAuthValue+jwt)

	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Result().StatusCode)

	type response struct {
		ID           string    `json:"id"`
		UserID       string    `json:"user_id"`
		SourceID     string    `json:"source_id"`
		TargetID     string    `json:"target_id"`
		SourceFormat string    `json:"source_format"`
		TargetFormat string    `json:"target_format"`
		Ratio        int       `json:"ratio"`
		Created      time.Time `json:"created"`
		Updated      time.Time `json:"updated"`
		Status       string    `json:"status"`
	}
	var history []response

	err = json.Unmarshal(w.Body.Bytes(), &history)
	s.NoError(err)

	s.Equal(userID, history[0].UserID)
	s.Equal(requestID, history[0].ID)
	s.Equal(defaultSourceFormat, history[0].SourceFormat)
	s.Equal(defaultTargetFormat, history[0].TargetFormat)
}

func (s *APITestSuite) TestDownloadImage() {
	userID, err := s.repos.users.InsertUser(context.Background(), defaultEmail, defaultPassword)
	s.NoError(err)

	sourceID, err := s.repos.images.InsertImage(context.Background(), defaultFile, defaultSourceFormat)
	s.NoError(err)

	_, err = s.repos.requests.InsertRequest(context.Background(), userID, sourceID, defaultSourceFormat, defaultTargetFormat, 99)
	s.NoError(err)

	jwt, err := s.tm.GenerateAccessToken(userID)
	s.NoError(err)

	s.mocks.storageMock.EXPECT().GetDownloadURL(sourceID).Return(defaultImageURL, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, testDownloadURL, nil)
	req.Header.Set(headerTypeKey, headerTypeValue)
	req.Header.Set(headerAuthKey, headerAuthValue+jwt)
	q := req.URL.Query()
	q.Add("id", sourceID)
	req.URL.RawQuery = q.Encode()

	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Result().StatusCode)

	type response struct {
		ImageURL string `json:"image_url"`
	}
	var resp response

	err = json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)

	s.Equal(resp.ImageURL, defaultImageURL)
}

func createMockRequest(t *testing.T, url, method string) *http.Request {
	file, err := os.Create(defaultFilename)
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
	part, err := writer.CreateFormFile("file", filepath.Base(defaultFilename))
	require.NoError(t, err)
	_, _ = io.Copy(part, file)

	writer.WriteField("targetFormat", defaultTargetFormat)
	writer.WriteField("ratio", defaultRatio)

	err = writer.Close()
	require.NoError(t, err)

	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)

	req.Header.Set(headerTypeKey, writer.FormDataContentType())

	return req
}

func (s *APITestSuite) TestConvertImage() {
	userID, err := s.repos.users.InsertUser(context.Background(), defaultEmail, defaultPassword)
	s.NoError(err)

	jwt, err := s.tm.GenerateAccessToken(userID)
	s.NoError(err)

	s.mocks.storageMock.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(nil)
	s.mocks.producerMock.EXPECT().
		SendToQueue(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	w := httptest.NewRecorder()
	req := createMockRequest(s.T(), testConversionURL, http.MethodPost)
	defer os.Remove(defaultFilename)
	req.Header.Set(headerAuthKey, headerAuthValue+jwt)

	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusAccepted, w.Result().StatusCode)
}
