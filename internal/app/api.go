// Package app provides a function to start the application.
package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Konstantsiy/image-converter/internal/service"

	"github.com/Konstantsiy/image-converter/pkg/logger"

	"github.com/Konstantsiy/image-converter/internal/queue"

	"github.com/Konstantsiy/image-converter/internal/config"
	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/internal/server"
	"github.com/Konstantsiy/image-converter/internal/storage"
	"github.com/Konstantsiy/image-converter/pkg/jwt"
	"github.com/gorilla/mux"
)

// Start starts the application server.
func Start() error {
	r := mux.NewRouter()

	conf, err := config.Load()
	if err != nil {
		return fmt.Errorf("can't load configs: %w", err)
	}

	db, err := repository.NewPostgresDB(conf.DBConf)
	if err != nil {
		return fmt.Errorf("can't connect to postgres database: %v", err)
	}
	defer db.Close()

	logger.FromContext(context.Background()).
		WithField("host", conf.DBConf.Host).WithField("port", conf.DBConf.Port).
		Infoln("database connected successfully")

	tokenManager, err := jwt.NewTokenManager(conf.JWTConf)
	if err != nil {
		return fmt.Errorf("token manager error: %w", err)
	}
	logger.FromContext(context.Background()).Infoln("JWT-manager created successfully")

	st, err := storage.NewStorage(conf.AWSConf)
	if err != nil {
		return fmt.Errorf("can't create storage: %v", err)
	}
	logger.FromContext(context.Background()).Infoln("AWS S3 connected successfully")

	producer, err := queue.NewRabbitMQProducer(conf.RabbitMQConf)
	if err != nil {
		return fmt.Errorf("can't create producer: %w", err)
	}
	logger.FromContext(context.Background()).Infoln("RabbitMQ client (producer) initialized successfully")

	usersRepo, err := repository.NewUsersRepository(db)
	if err != nil {
		return fmt.Errorf("users repository creating error: %w", err)
	}

	imageRepo, err := repository.NewImagesRepository(db)
	if err != nil {
		return fmt.Errorf("images repository creating error: %w", err)
	}

	requestsRepo, err := repository.NewRequestsRepository(db)
	if err != nil {
		return fmt.Errorf("requests repository creating error: %w", err)
	}

	authService := service.NewAuthService(usersRepo, tokenManager)
	imagesService := service.NewImageService(imageRepo, requestsRepo, st)
	requestsService := service.NewRequestsService(requestsRepo)

	s := server.NewServer(authService, imagesService, requestsService, producer)
	s.RegisterRoutes(r)

	return http.ListenAndServe(":"+conf.AppPort, r)
}
