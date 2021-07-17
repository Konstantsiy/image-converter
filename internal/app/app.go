// Package app provides a function to start the application.
package app

import (
	"net/http"

	"github.com/Konstantsiy/image-converter/internal/auth"
	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/internal/server"
	"github.com/gorilla/mux"
)

// Start starts the application server.
func Start() error {
	r := mux.NewRouter()
	repo := repository.NewRepository()
	tokenManager := auth.NewTokenManager("NHbNaO6LFERWnOUbU7l3MJdmCailwSzjO76O")
	s := server.NewServer(repo, tokenManager)
	s.RegisterRoutes(r)
	return http.ListenAndServe(":8080", r)
}
