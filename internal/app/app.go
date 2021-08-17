// Package app provides a function to start the application.
package app

import (
	"fmt"
	"net/http"

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

	db, err := repository.NewPostgresDB(&conf)
	if err != nil {
		return fmt.Errorf("can't connect to postgres database: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	tokenManager := auth.NewTokenManager(conf.PublicKey, conf.PrivateKey)

	s := server.NewServer(repo, tokenManager)
	s.RegisterRoutes(r)
	return http.ListenAndServe(":8080", r)
}
