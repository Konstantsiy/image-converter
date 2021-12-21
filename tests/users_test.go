package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Konstantsiy/image-converter/pkg/hash"

	"github.com/Konstantsiy/image-converter/internal/repository"
)

func (s *APITestSuite) truncateTableUsers(db *sql.DB) {
	query := "TRUNCATE TABLE converter.users CASCADE;"
	_, err := db.Exec(query)
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
	s.Equal(http.StatusOK, w.Result().StatusCode)

	type response struct {
		UserID string `json:"user_id"`
	}
	var resp response
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	s.NoError(err)

	ur, err := repository.NewUsersRepository(s.db)
	s.NoError(err)

	user, err := ur.GetUserByEmail(context.Background(), defaultEmail)
	s.NoError(err)

	s.Equal(defaultEmail, user.Email)
	equal, err := hash.ComparePasswordHash(defaultPassword, user.Password)
	s.NoError(err)
	s.True(equal)
}
