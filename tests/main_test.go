// Package test contains integration tests.
package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/suite"

	"github.com/Konstantsiy/image-converter/internal/server"

	"github.com/Konstantsiy/image-converter/internal/config"
	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/internal/service"
	"github.com/Konstantsiy/image-converter/pkg/jwt"
)

type APITestSuite struct {
	suite.Suite

	db     *sql.DB
	conf   *config.Config
	serv   *server.Server
	router *mux.Router
	tm     *jwt.TokenManager
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

func (s *APITestSuite) FailWithError(err error) {
	s.FailNow(err.Error())
}

func (s *APITestSuite) SetupSuite() {
	conf, err := config.Load()
	s.conf = &conf
	s.NoError(err)
	if err != nil {
		s.FailWithError(fmt.Errorf("can't load configs: %w", err))
	}
	s.initDependencies(&conf)
}

func (s *APITestSuite) TearDownSuite() {
	err := s.db.Close()
	if err != nil {
		s.FailWithError(err)
	}
}

func (s *APITestSuite) SetupTest() {
	s.truncateTableUsers()
}

func (s *APITestSuite) TearDownTest() {
	s.truncateTableUsers()
}

func (s *APITestSuite) initDependencies(conf *config.Config) {
	db, err := repository.NewPostgresDB(conf.DBConf)
	if err != nil {
		s.FailWithError(fmt.Errorf("can't connect to postgres database: %v", err))
	}
	s.db = db

	tm, err := jwt.NewTokenManager(conf.JWTConf)
	if err != nil {
		s.FailWithError(fmt.Errorf("token manager error: %w", err))
	}
	s.tm = tm

	usersRepo, err := repository.NewUsersRepository(s.db)
	if err != nil {
		s.FailWithError(fmt.Errorf("users repository creating error: %w", err))
	}

	imageRepo, err := repository.NewImagesRepository(s.db)
	if err != nil {
		s.FailWithError(fmt.Errorf("images repository creating error: %w", err))
	}

	requestsRepo, err := repository.NewRequestsRepository(s.db)
	if err != nil {
		s.FailWithError(fmt.Errorf("requests repository creating error: %w", err))
	}

	authService := service.NewAuthService(usersRepo, s.tm)
	imagesService := service.NewImageService(imageRepo, requestsRepo, nil, nil)
	requestsService := service.NewRequestsService(requestsRepo)

	s.serv = server.NewServer(authService, imagesService, requestsService, nil)
	s.router = mux.NewRouter()
	s.serv.RegisterRoutes(s.router)
}
