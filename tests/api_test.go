package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/pkg/hash"
)

func (s *APITestSuite) truncateTableUsers() {
	query := "TRUNCATE TABLE converter.users CASCADE;"
	_, err := s.db.Exec(query)
	if err != nil {
		s.FailWithError(fmt.Errorf("unable to truncate users table: %v", err))
	}
}

func (s *APITestSuite) TestUserSignIn() {
	const (
		defaultEmail    = "email1@gmail.com"
		defaultPassword = "Password1"
		testURL         = "/user/login"
		headerTypeKey   = "Content-type"
		headerTypeValue = "application/json"
	)

	signInData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, defaultEmail, defaultPassword)
	pwdHash, err := hash.GeneratePasswordHash(defaultPassword)
	s.NoError(err)

	ur, err := repository.NewUsersRepository(s.db)
	s.NoError(err)

	_, err = ur.InsertUser(context.Background(), defaultEmail, pwdHash)
	s.NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, testURL, bytes.NewBuffer([]byte(signInData)))
	req.Header.Set(headerTypeKey, headerTypeValue)

	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.truncateTableUsers()
}

func (s *APITestSuite) TestUserSignUp() {
	const (
		defaultEmail    = "email2@gmail.com"
		defaultPassword = "Password223"
		testURL         = "/user/signup"
		headerTypeKey   = "Content-type"
		headerTypeValue = "application/json"
	)

	signUpData := fmt.Sprintf(`{"email" :"%s","password": "%s"}`, defaultEmail, defaultPassword)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, testURL, bytes.NewBuffer([]byte(signUpData)))
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
	const (
		defaultEmail        = "email3@gmail.com"
		defaultPassword     = "Password3"
		testURL             = "/requests"
		headerTypeKey       = "Content-type"
		headerTypeValue     = "application/json"
		headerAuthKey       = "Authorization"
		headerAuthValue     = "Bearer "
		defaultFileName     = "file1"
		defaultSourceFormat = "jpg"
		defaultTargetFormat = "png"
	)

	ur, err := repository.NewUsersRepository(s.db)
	s.NoError(err)
	userID, err := ur.InsertUser(context.Background(), defaultEmail, defaultPassword)
	s.NoError(err)

	ir, err := repository.NewImagesRepository(s.db)
	s.NoError(err)
	sourceID, err := ir.InsertImage(context.Background(), defaultFileName, defaultSourceFormat)
	s.NoError(err)

	rr, err := repository.NewRequestsRepository(s.db)
	s.NoError(err)
	requestID, err := rr.InsertRequest(context.Background(), userID, sourceID, defaultSourceFormat, defaultTargetFormat, 99)
	s.NoError(err)

	jwt, err := s.tm.GenerateAccessToken(userID)
	s.NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
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
