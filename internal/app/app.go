// Package app provides a function to start the application.
package app

import (
	"net/http"

	"github.com/Konstantsiy/image-converter/internal/converter"

	"github.com/Konstantsiy/image-converter/internal/auth"
	"github.com/Konstantsiy/image-converter/internal/config"
	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/internal/server"
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
	conv := converter.NewConverter()

	s := server.NewServer(repo, tokenManager, conv)
	s.RegisterRoutes(r)
	return http.ListenAndServe(":8080", r)
}
