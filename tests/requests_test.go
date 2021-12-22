package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Konstantsiy/image-converter/internal/repository"
)

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

	s.truncateTableUsers()

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
