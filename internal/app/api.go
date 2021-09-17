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

	logger.Info(context.Background(), "port: "+conf.AppPort)

	db, err := repository.NewPostgresDB(&conf)
	if err != nil {
		return fmt.Errorf("can't connect to postgres database: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	tokenManager, err := jwt.NewTokenManager(&conf)
	if err != nil {
		return fmt.Errorf("token manager error: %w", err)
	}

	st, err := storage.NewStorage(storage.S3Config{
		Region:          conf.Region,
		AccessKeyID:     conf.AccessKeyID,
		SecretAccessKey: conf.SecretAccessKey,
		BucketName:      conf.BucketName,
	})
	if err != nil {
		return fmt.Errorf("can't create storage: %v", err)
	}

	producer, err := queue.NewProducer(&conf)
	if err != nil {
		return fmt.Errorf("can't create publisher: %w", err)
	}

	s := server.NewServer(repo, tokenManager, st, producer)
	s.RegisterRoutes(r)
	return http.ListenAndServe(":"+conf.AppPort, r)
}
