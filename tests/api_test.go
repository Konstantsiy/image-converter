// Package test contains integration tests.
package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Konstantsiy/image-converter/pkg/hash"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"

	"github.com/Konstantsiy/image-converter/internal/server"

	"github.com/Konstantsiy/image-converter/internal/config"
	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/internal/service"
	"github.com/Konstantsiy/image-converter/pkg/jwt"
)

func initServer(t *testing.T, conf *config.Config, db *sql.DB) *server.Server {
	tm, err := jwt.NewTokenManager(conf.JWTConf)
	if err != nil {
		t.Fatalf("token manager error: %s", err.Error())
	}

	usersRepo, err := repository.NewUsersRepository(db)
	if err != nil {
		t.Fatalf("users repository creating error: %s", err.Error())
	}

	imageRepo, err := repository.NewImagesRepository(db)
	if err != nil {
		t.Fatalf("images repository creating error: %s", err.Error())
	}

	requestsRepo, err := repository.NewRequestsRepository(db)
	if err != nil {
		t.Fatalf("requests repository creating error: %s", err.Error())
	}

	authService := service.NewAuthService(usersRepo, tm)
	imagesService := service.NewImageService(imageRepo, requestsRepo, nil, nil)
	requestsService := service.NewRequestsService(requestsRepo)

	return server.NewServer(authService, imagesService, requestsService, nil)
}

func truncateTableUsers(t *testing.T, db *sql.DB) {
	query := "TRUNCATE TABLE converter.users CASCADE;"
	_, err := db.Exec(query)
	if err != nil {
		t.Error(fmt.Errorf("unable to truncate users table: %v", err))
	}
}

func TestUsersSignUp(t *testing.T) {
	const (
		driverName      = "postgres"
		defaultEmail    = "email2@gmail.com"
		defaultPassword = "Password223"
		testURL         = "/user/signup"
		headerTypeKey   = "Content-type"
		headerTypeValue = "application/json"
	)

	conf, err := config.Load()
	assert.NoError(t, err)

	db, err := sql.Open(
		driverName,
		fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", conf.DBConf.User, conf.DBConf.Password,
			conf.DBConf.Host, conf.DBConf.Port, conf.DBConf.DBName, conf.DBConf.SSLMode))
	assert.NoError(t, err)

	err = db.Ping()
	assert.NoError(t, err)

	truncateTableUsers(t, db)

	signUpData := fmt.Sprintf(`{"email" :"%s","password": "%s"}`, defaultEmail, defaultPassword)

	s := initServer(t, &conf, db)

	r := mux.NewRouter()
	s.RegisterRoutes(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, testURL, bytes.NewBuffer([]byte(signUpData)))
	req.Header.Set(headerTypeKey, headerTypeValue)

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	type response struct {
		UserID string `json:"user_id"`
	}
	var resp response
	err = json.Unmarshal([]byte(w.Body.String()), &resp)
	assert.NoError(t, err)

	ur, err := repository.NewUsersRepository(db)
	assert.NoError(t, err)

	user, err := ur.GetUserByEmail(context.Background(), defaultEmail)
	assert.NoError(t, err)

	assert.Equal(t, defaultEmail, user.Email)
	equal, err := hash.ComparePasswordHash(defaultPassword, user.Password)
	assert.NoError(t, err)
	assert.True(t, equal)
}

//type APITestSuite struct {
//	suite.Suite
//
//	db *sql.DB
//
//	authService     *service.AuthService
//	imagesService   *service.ImageService
//	requestsService *service.RequestsService
//
//	tm *jwt.TokenManager
//}
//
//func TestAPISuite(t *testing.T) {
//	suite.Run(t, new(APITestSuite))
//}
//
//func TestMain(m *testing.M) {
//	rc := m.Run()
//	os.Exit(rc)
//}
//
//func (s *APITestSuite) FailWithError(err error) {
//	s.FailNow(err.Error())
//}
//
//func (s *APITestSuite) SetupSuite() {
//	conf, err := config.Load()
//	if err != nil {
//		s.FailWithError(fmt.Errorf("can't load configs: %w", err))
//	}
//
//	s.initDependencies(&conf)
//}
//
//func (s *APITestSuite) TearDownSuite() {
//	err := s.db.Close()
//	if err != nil {
//		s.FailWithError(err)
//	}
//}
//
//func (s *APITestSuite) initDependencies(conf *config.Config) {
//	db, err := repository.NewPostgresDB(conf.DBConf)
//	if err != nil {
//		s.FailWithError(fmt.Errorf("can't connect to postgres database: %v", err))
//	}
//	s.db = db
//
//	tm, err := jwt.NewTokenManager(conf.JWTConf)
//	if err != nil {
//		s.FailWithError(fmt.Errorf("token manager error: %w", err))
//	}
//	s.tm = tm
//
//	usersRepo, err := repository.NewUsersRepository(s.db)
//	if err != nil {
//		s.FailWithError(fmt.Errorf("users repository creating error: %w", err))
//	}
//
//	imageRepo, err := repository.NewImagesRepository(s.db)
//	if err != nil {
//		s.FailWithError(fmt.Errorf("images repository creating error: %w", err))
//	}
//
//	requestsRepo, err := repository.NewRequestsRepository(s.db)
//	if err != nil {
//		s.FailWithError(fmt.Errorf("requests repository creating error: %w", err))
//	}
//
//	s.authService = service.NewAuthService(usersRepo, s.tm)
//	s.imagesService = service.NewImageService(imageRepo, requestsRepo, nil, nil)
//	s.requestsService = service.NewRequestsService(requestsRepo)
//}
