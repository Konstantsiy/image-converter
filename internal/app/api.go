// Package app provides a function to start the application.
package app

import (
	"context"
	"fmt"
	"net/http"

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

	logger.Info(context.Background(), fmt.Sprintf("configs:\nport: %s\nDB: %+v\nAWS: %+v\nJWT: %+v\nRabbit: %+v\n",
		conf.AppPort, conf.DBConf, conf.AWSConf, conf.JWTConf, conf.RabbitMQConf))

	db, err := repository.NewPostgresDB(conf.DBConf)
	if err != nil {
		return fmt.Errorf("can't connect to postgres database: %v", err)
	}
	defer db.Close()
	logger.Info(context.Background(), fmt.Sprintf("database connected successfully (%s:%s)",
		conf.DBConf.Host, conf.DBConf.Port))

	repo := repository.NewRepository(db)

	tokenManager, err := jwt.NewTokenManager(conf.JWTConf)
	if err != nil {
		return fmt.Errorf("token manager error: %w", err)
	}
	logger.Info(context.Background(), "JWT-manager created successfully")

	st, err := storage.NewStorage(conf.AWSConf)
	if err != nil {
		return fmt.Errorf("can't create storage: %v", err)
	}
	logger.Info(context.Background(), "AWS S3 connected successfully")

	producer, err := queue.NewProducer(conf.RabbitMQConf)
	if err != nil {
		return fmt.Errorf("can't create producer: %w", err)
	}
	logger.Info(context.Background(), "RabbitMQ client (producer) initialized successfully")

	s := server.NewServer(repo, tokenManager, st, producer)
	s.RegisterRoutes(r)

	return http.ListenAndServe(":"+conf.AppPort, r)
}
