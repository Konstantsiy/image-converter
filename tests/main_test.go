// Package test contains integration tests.
package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	mockqueue "github.com/Konstantsiy/image-converter/internal/queue/mock"

	"github.com/golang/mock/gomock"

	mockstorage "github.com/Konstantsiy/image-converter/internal/storage/mock"

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

	repos *repositories
	mocks *mocks
	mc    *gomock.Controller
}

type repositories struct {
	images   repository.Images
	requests repository.Requests
	users    repository.Users
}

type mocks struct {
	storageMock  *mockstorage.MockStorage
	producerMock *mockqueue.MockProducer
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
	s.T().Log("setup test suite configuration")
	conf, err := config.Load()
	s.conf = &conf
	s.NoError(err)
	if err != nil {
		s.FailWithError(fmt.Errorf("can't load configs: %w", err))
	}
	s.initMocks()
	s.initDependencies(&conf)
}

func (s *APITestSuite) TearDownSuite() {
	s.T().Log("tear down test suite: finish mock controller & close database connection")
	s.mc.Finish()
	err := s.db.Close()
	if err != nil {
		s.FailWithError(err)
	}
}

func (s *APITestSuite) SetupTest() {
	s.T().Log("setup test: truncate tables in database")
	s.truncateTables()
}

func (s *APITestSuite) TearDownTest() {
	s.T().Log("tear down test: truncate tables in database")
	s.truncateTables()
}

func (s *APITestSuite) initMocks() {
	s.T().Log("init mocks")
	s.mc = gomock.NewController(s.T())
	s.mocks = &mocks{
		storageMock:  mockstorage.NewMockStorage(s.mc),
		producerMock: mockqueue.NewMockProducer(s.mc),
	}
}

func (s *APITestSuite) initDependencies(conf *config.Config) {
	s.T().Log("open database connection")
	db, err := repository.NewPostgresDB(conf.DBConf)
	if err != nil {
		s.FailWithError(fmt.Errorf("can't connect to postgres database: %v", err))
	}
	s.db = db

	s.T().Log("init token manager")
	tm, err := jwt.NewTokenManager(conf.JWTConf)
	if err != nil {
		s.FailWithError(fmt.Errorf("token manager error: %w", err))
	}
	s.tm = tm

	s.T().Log("init repositories")
	usersRepo, err := repository.NewUsersRepository(s.db)
	if err != nil {
		s.FailWithError(fmt.Errorf("users repository creating error: %w", err))
	}

	imagesRepo, err := repository.NewImagesRepository(s.db)
	if err != nil {
		s.FailWithError(fmt.Errorf("images repository creating error: %w", err))
	}

	requestsRepo, err := repository.NewRequestsRepository(s.db)
	if err != nil {
		s.FailWithError(fmt.Errorf("requests repository creating error: %w", err))
	}

	s.repos = &repositories{
		users:    usersRepo,
		images:   imagesRepo,
		requests: requestsRepo,
	}

	authService := service.NewAuthService(usersRepo, s.tm)
	imagesService := service.NewImageService(imagesRepo, requestsRepo, s.mocks.storageMock)
	requestsService := service.NewRequestsService(requestsRepo)

	s.T().Log("init application server")
	s.serv = server.NewServer(authService, imagesService, requestsService, s.mocks.producerMock)
	s.router = mux.NewRouter()
	s.T().Log("register http routing")
	s.serv.RegisterRoutes(s.router)
}

func (s *APITestSuite) truncateTables() {
	query1 := "TRUNCATE TABLE converter.users CASCADE;"
	query2 := "TRUNCATE TABLE converter.images CASCADE;"
	query3 := "TRUNCATE TABLE converter.requests;"

	_, err := s.db.Exec(query1)
	if err != nil {
		s.FailWithError(fmt.Errorf("unable to truncate users table: %v", err))
	}
	_, err = s.db.Exec(query2)
	if err != nil {
		s.FailWithError(fmt.Errorf("unable to truncate images table: %v", err))
	}

	_, err = s.db.Exec(query3)
	if err != nil {
		s.FailWithError(fmt.Errorf("unable to truncate requests table: %v", err))
	}
}
