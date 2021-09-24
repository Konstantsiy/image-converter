// Package server implements http handlers.
package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Konstantsiy/image-converter/pkg/logger"

	"github.com/Konstantsiy/image-converter/internal/queue"

	"github.com/Konstantsiy/image-converter/internal/appcontext"
	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/internal/storage"
	"github.com/Konstantsiy/image-converter/internal/validation"
	"github.com/Konstantsiy/image-converter/pkg/hash"
	"github.com/Konstantsiy/image-converter/pkg/jwt"
	"github.com/gorilla/mux"
)

const (
	AuthorizationHeader  = "Authorization"
	NeededSecurityScheme = "Bearer"

	DefaultStatusCode = 200
)

// AuthRequest represents the user's authorization request.
type AuthRequest struct {
	Email    string
	Password string
}

// ConvertRequest represents an image conversion request.
type ConvertRequest struct {
	File         string
	SourceFormat string
	TargetFormat string
	Ratio        int
}

// ConvertResponse represents an image conversion response.
type ConvertResponse struct {
	RequestID string `json:"request_id"`
}

// LoginResponse represents token for authorization response.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// SignUpResponse represents user id from sign up response.
type SignUpResponse struct {
	UserID string `json:"user_id"`
}

// DownloadResponse represents downloaded image URL.
type DownloadResponse struct {
	ImageURL string `json:"image_url"`
}

// Server represents application server.
type Server struct {
	repo         *repository.Repository
	tokenManager *jwt.TokenManager
	storage      *storage.Storage
	producer     *queue.Producer
}

// NewServer creates new application server.
func NewServer(repo *repository.Repository, tokenManager *jwt.TokenManager, storage *storage.Storage, producer *queue.Producer) *Server {
	return &Server{
		repo:         repo,
		tokenManager: tokenManager,
		storage:      storage,
		producer:     producer,
	}
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
	var request AuthRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()

	user, err := s.repo.GetUserByEmail(request.Email)
	if err == repository.ErrNoSuchUser {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ok, err := hash.ComparePasswords(user.Password, request.Password); !ok || err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID)
	if err != nil {
		http.Error(w, "can't generate access token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	refreshToken, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "can't generate refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(w, LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, http.StatusOK)
}

// SignUp implements the user registration process.
func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	var request AuthRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()

	if err = validation.ValidateSignUpRequest(request.Email, request.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashPwd, err := hash.GeneratePasswordHash(request.Password)
	if err != nil {
		http.Error(w, "can't generate password hash: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userID, err := s.repo.InsertUser(request.Email, hashPwd)
	if errors.Is(err, repository.ErrUserAlreadyExists) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(w, SignUpResponse{UserID: userID}, http.StatusOK)
}

// ConvertImage converts needed image according to the request.
func (s *Server) ConvertImage(w http.ResponseWriter, r *http.Request) {
	sourceFile, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "can't get sourceFile from form", http.StatusBadRequest)
		return
	}
	defer sourceFile.Close()

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

	userID, ok := appcontext.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "can't get user id from application context", http.StatusInternalServerError)
		return
	}

	sourceFileID, err := s.repo.InsertImage(filename, sourceFormat)
	if err != nil {
		http.Error(w, fmt.Sprintf("repository error: %v", err), http.StatusInternalServerError)
		return
	}

	err = s.storage.UploadFile(sourceFile, sourceFileID)
	if err != nil {
		http.Error(w, fmt.Sprintf("storage error: %v", err), http.StatusInternalServerError)
		return
	}

	requestID, err := s.repo.MakeRequest(userID, sourceFileID, sourceFormat, targetFormat, ratio)
	if err != nil {
		http.Error(w, fmt.Sprintf("repository error: %v", err), http.StatusInternalServerError)
		return
	}

	sendResponse(w, ConvertResponse{RequestID: requestID}, http.StatusAccepted)

	err = s.producer.SendToQueue(sourceFileID, filename, sourceFormat, targetFormat, requestID, ratio)
	if err != nil {
		http.Error(w, fmt.Sprint("can't send data to queue: %w", err), http.StatusInternalServerError)
		return
	}
}

// DownloadImage allows you to download original/converted image by id.
func (s *Server) DownloadImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	imageID, err := s.repo.GetImageIDInStore(id)
	if errors.Is(err, repository.ErrNoSuchImage) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url, err := s.storage.GetDownloadURL(imageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(w, DownloadResponse{ImageURL: url}, http.StatusOK)
}

// GetRequestsHistory displays the user's request history.
func (s *Server) GetRequestsHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := appcontext.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "can't get user id from application context", http.StatusInternalServerError)
		return
	}

	requestsHistory, err := s.repo.GetRequestsByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(w, requestsHistory, http.StatusOK)
}

// sendResponse marshals and writes response to the ResponseWriter.
func sendResponse(w http.ResponseWriter, resp interface{}, code int) {
	respJSON, err := json.Marshal(resp)
	if err != nil {
		logger.Error(context.Background(), fmt.Errorf("can't marshal response: %v", err))
		fmt.Fprint(w, resp)
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}
