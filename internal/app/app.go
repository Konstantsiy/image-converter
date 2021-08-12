// Package app provides a function to start the application.
package app

import (
	"fmt"
	"net/http"

	"github.com/Konstantsiy/image-converter/internal/auth"
	"github.com/Konstantsiy/image-converter/internal/config"
	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/internal/server"
	"github.com/Konstantsiy/image-converter/internal/storage"
	"github.com/gorilla/mux"
)

// Start starts the application server.
func Start() error {
	r := mux.NewRouter()

	conf, err := config.Load()
	if err != nil {
		return err
	}

	repo := repository.NewRepository()
	tokenManager := auth.NewTokenManager(conf.PublicKey, conf.PrivateKey)

	st, err := storage.NewStorage(storage.S3Config{
		Region:          conf.Region,
		AccessKeyID:     conf.AccessKeyID,
		SecretAccessKey: conf.SecretAccessKey,
		BucketName:      conf.BucketName,
	})
	if err != nil {
		return fmt.Errorf("storage error: %v", err)
	}

	s := server.NewServer(repo, tokenManager, st)
	s.RegisterRoutes(r)
	return http.ListenAndServe(":"+conf.Port, r)
}
