// Package server implements http handlers.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Konstantsiy/image-converter/internal/auth"
	"github.com/Konstantsiy/image-converter/internal/domain"
	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/internal/validation"
	"github.com/gorilla/mux"
)

// Server represents application server.
type Server struct {
	repo         *repository.Repository
	tokenManager *auth.TokenManager
}

// NewServer creates new application server.
func NewServer(repo *repository.Repository, tokenManager *auth.TokenManager) *Server {
	return &Server{
		repo:         repo,
		tokenManager: tokenManager,
	}
}

// AuthMiddleware checks user authorization.
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check JWT tokens
		next.ServeHTTP(w, r)
	})
}

// RegisterRoutes registers application routers.
func (s *Server) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/user/login", s.LogIn).Methods("POST")
	r.HandleFunc("/user/signup", s.SignUp).Methods("POST")

	api := r.NewRoute().Subrouter()
	api.Use(s.AuthMiddleware)

	api.HandleFunc("/conversion", s.ConvertImage).Methods("POST")
	api.HandleFunc("/images/{id}", s.DownloadImage).Methods("GET")
	api.HandleFunc("/requests", s.GetRequestsHistory).Methods("GET")
}

// LogIn implements the user authentication process.
func (s *Server) LogIn(w http.ResponseWriter, r *http.Request) {
	var request domain.AuthRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()

	var user domain.UserInfo
	user, err = s.repo.GetUserByEmail(request.Email)
	if err == repository.ErrNoSuchUser {
		http.Error(w, "can't get user info: "+err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tokens domain.TokensResponse
	tokens.AccessToken, err = s.tokenManager.GenerateAccessToken(user.ID)
	if err != nil {
		http.Error(w, "can't generate access token: "+err.Error(), http.StatusUnauthorized)
		return
	}
	tokens.RefreshToken, err = s.tokenManager.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "can't generate refresh token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	fmt.Fprint(w, tokens.AccessToken, tokens.RefreshToken)
}

// SignUp implements the user registration process.
func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	var request domain.AuthRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()

	if err = validation.ValidateUserCredentials(request.Email, request.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := s.repo.InsertUser(request.Email, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	fmt.Fprint(w, userID)
}

// ConvertImage converts needed image according to the request.
func (s *Server) ConvertImage(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "can't get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	sourceFormat := r.FormValue("sourceFormat")
	targetFormat := r.FormValue("targetFormat")
	filename := strings.TrimSuffix(header.Filename, "."+sourceFormat)
	ratio, err := strconv.Atoi(r.FormValue("ratio"))
	if err != nil {
		http.Error(w, "invalid ratio form value", http.StatusBadRequest)
		return
	}

	if err = validation.ValidateConversionRequest(filename, sourceFormat, targetFormat, ratio); err != nil {
		http.Error(w, fmt.Sprint(err.Error()), http.StatusBadRequest)
		return
	}

	// ...
}

// DownloadImage allows you to download original/converted image by id.
func (s *Server) DownloadImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := vars["id"]

	image, err := s.repo.GetImageByID(imageID)
	if err == repository.ErrNoSuchImage {
		http.Error(w, "can't get image info: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get image downloaded URL from storage by image.Location

	url := "http(s)://s3.amazonaws.com/" + image.Location + "/file_name.extension"

	fmt.Fprint(w, url)
}

// GetRequestsHistory displays the user's request history.
func (s *Server) GetRequestsHistory(w http.ResponseWriter, r *http.Request) {
	// get userID from application context?

	userID := "7186afcc-cae7-11eb-80ff-0bc45a674b3c"
	requestsHistory, err := s.repo.GetRequestsByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, fmt.Sprint(requestsHistory))
}
