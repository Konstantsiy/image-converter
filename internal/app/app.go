// Package app provides a function to start the application.
package app

import (
	"fmt"
	"net/http"

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
		return err
	}

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
		return fmt.Errorf("storage error: %v", err)
	}

	producer, err := queue.NewProducer(&conf)
	if err != nil {
		return fmt.Errorf("can't create publisher: %w", err)
	}

	consumer, err := queue.NewConsumer(repo, st, &conf)
	if err != nil {
		return fmt.Errorf("can't create consumer: %w", err)
	}

	err = consumer.ListenToQueue()
	if err != nil {
		return fmt.Errorf("queue consumer error: %w", err)
	}

	s := server.NewServer(repo, tokenManager, st, producer)
	s.RegisterRoutes(r)
	return http.ListenAndServe(":8080", r)
}
