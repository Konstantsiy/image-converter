// Package app provides a function to start the application.
package app

import (
	"net/http"

	"github.com/Konstantsiy/image-converter/internal/server"
	"github.com/gorilla/mux"
)

// Start starts the application server.
func Start() error {
	r := mux.NewRouter()
	s := server.NewServer()
	s.RegisterRoutes(r)
	return http.ListenAndServe(":8080", r)
}
