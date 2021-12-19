// Package test contains integration tests.
package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	mockqueue "github.com/Konstantsiy/image-converter/internal/queue/mock"
	mockstorage "github.com/Konstantsiy/image-converter/internal/storage/mock"

	"github.com/Konstantsiy/image-converter/internal/config"
	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/internal/service"
	"github.com/Konstantsiy/image-converter/pkg/jwt"
	"github.com/stretchr/testify/suite"
)

type mocks struct {
	StorageMock  *mockstorage.MockStorage
	ProducerMock *mockqueue.MockProducer
}

type APITestSuite struct {
	suite.Suite

	db *sql.DB

	authService     *service.AuthService
	imagesService   *service.ImageService
	requestsService *service.RequestsService

	tm    *jwt.TokenManager
	mocks *mocks
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
	if err != nil {
		s.FailWithError(fmt.Errorf("can't load configs: %w", err))
	}

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

	s.initMocks()
	s.initDependencies()
}

func (s *APITestSuite) initDependencies() {
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

	s.authService = service.NewAuthService(usersRepo, s.tm)
	s.imagesService = service.NewImageService(imageRepo, requestsRepo, nil, nil)
	s.requestsService = service.NewRequestsService(requestsRepo)
}

func (s *APITestSuite) initMocks() {}
