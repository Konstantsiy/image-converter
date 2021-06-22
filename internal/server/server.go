// Package server implements http handlers.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Konstantsiy/image-converter/internal/model"
	"github.com/Konstantsiy/image-converter/internal/service"
	"github.com/gorilla/mux"
)

// Server represents application server.
type Server struct {
	userService    service.UserService
	imageService   service.ImageService
	requestService service.RequestService
}

// NewServer creates new application server.
func NewServer() *Server {
	return &Server{
		userService:    service.UserService{},
		imageService:   service.ImageService{},
		requestService: service.RequestService{},
	}
}

// RegisterRoutes registers application routers.
func (s *Server) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/user/login", s.LogIn).Methods("POST")
	r.HandleFunc("/user/signup", s.SignUp).Methods("POST")
	r.HandleFunc("/conversion", s.ConvertImage).Methods("POST")
	r.HandleFunc("/images/{id}", s.DownloadImage).Methods("GET")
	r.HandleFunc("/requests", s.GetRequestsHistory).Methods("GET")
}

// LogIn implements the user authentication process.
func (s *Server) LogIn(w http.ResponseWriter, r *http.Request) {
	// validation middleware (for email and password)

	var request model.AuthRequest
	var err error

	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var tokens model.Tokens

	if tokens, err = s.userService.LogIn(request); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	fmt.Fprint(w, tokens.AccessToken)
}

// SignUp implements the user registration process.
func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	// validation middleware (for email and password)

	var request model.AuthRequest
	var err error

	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if s.userService.IsExists(request.Email) {
		http.Error(w, "a similar user is already registered in the system", http.StatusConflict)
		return
	}

	var userID int

	if userID, err = s.userService.SignUp(request); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	fmt.Fprint(w, userID)
}

// ConvertImage converts needed image according to the request.
func (s *Server) ConvertImage(w http.ResponseWriter, r *http.Request) {
	// auth middleware
	// validation middleware (for file, formats and ration)
	// upload image
	// convert image
	// work with storage
	// create a request
	// return request id
}

// DownloadImage allows you to download original/converted image by id.
func (s *Server) DownloadImage(w http.ResponseWriter, r *http.Request) {
	// auth middleware

	vars := mux.Vars(r)
	imageID := vars["id"]

	image, err := s.imageService.GetImageByID(imageID)
	if err != nil {
		http.Error(w, "can't get downloaded URL", http.StatusInternalServerError)
		return
	}

	url, err := s.imageService.GetDownloadedURLByLocation(image.Location)
	if err != nil {
		http.Error(w, "can't get downloaded URL", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, url)
}

// GetRequestsHistory displays the user's request history.
func (s *Server) GetRequestsHistory(w http.ResponseWriter, r *http.Request) {
	// auth middleware
	// get userID from application context?

	userID := "7186afcc-cae7-11eb-80ff-0bc45a674b3c"
	requestsHistory, err := s.requestService.GetRequestsHistory(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, fmt.Sprint(requestsHistory))
}
