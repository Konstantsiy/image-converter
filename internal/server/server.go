// Package server implements http handlers.
package server

import (
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

	IDQueryKey = "id"
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
	r.Use(s.LoggingMiddleware)
	r.HandleFunc("/user/login", s.LogIn).Methods("POST")
	r.HandleFunc("/user/signup", s.SignUp).Methods("POST")

	api := r.NewRoute().Subrouter()

	api.Use(s.AuthMiddleware)
	api.HandleFunc("/conversion", s.ConvertImage).Methods("POST")
	api.HandleFunc("/images", s.DownloadImage).Methods("GET")
	api.HandleFunc("/requests", s.GetRequestsHistory).Methods("GET")
}

// LogIn implements the user authentication process.
func (s *Server) LogIn(w http.ResponseWriter, r *http.Request) {
	var request AuthRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		reportError(w, err, http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()

	user, err := s.repo.GetUserByEmail(request.Email)
	if err == repository.ErrNoSuchUser {
		reportError(w, fmt.Errorf("invalid email or password"), http.StatusUnauthorized)
		return
	}
	if err != nil {
		reportError(w, err, http.StatusInternalServerError)
		return
	}

	if ok, err := hash.ComparePasswordHash(request.Password, user.Password); !ok || err != nil {
		reportError(w, fmt.Errorf("invalid email or password"), http.StatusUnauthorized)
		return
	}

	logger.FromContext(r.Context()).WithField("user_id", user.ID).Infoln("user successfully logged in")

	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID)
	if err != nil {
		reportError(w, fmt.Errorf("can't generate access token: %w", err), http.StatusInternalServerError)
		return
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		reportError(w, fmt.Errorf("can't generate refresh token: %w", err), http.StatusInternalServerError)
		return
	}

	sendResponse(w, LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, http.StatusOK)
}

// SignUp implements the user registration process.
func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	var request AuthRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		reportError(w, err, http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()

	if err = validation.ValidateSignUpRequest(request.Email, request.Password); err != nil {
		reportError(w, err, http.StatusBadRequest)
		return
	}

	hashPwd, err := hash.GeneratePasswordHash(request.Password)
	if err != nil {
		reportError(w, fmt.Errorf("can't generate password hash: %w", err), http.StatusInternalServerError)
		return
	}

	userID, err := s.repo.InsertUser(request.Email, hashPwd)
	if errors.Is(err, repository.ErrUserAlreadyExists) {
		reportError(w, err, http.StatusBadRequest)
		return
	}
	if err != nil {
		reportError(w, err, http.StatusInternalServerError)
		return
	}

	sendResponse(w, SignUpResponse{UserID: userID}, http.StatusOK)
}

// ConvertImage converts needed image according to the request.
func (s *Server) ConvertImage(w http.ResponseWriter, r *http.Request) {
	sourceFile, header, err := r.FormFile("file")
	if err != nil {
		reportError(w, fmt.Errorf("can't get sourceFile from form"), http.StatusBadRequest)
		return
	}
	defer sourceFile.Close()

	sourceFormat := r.FormValue("sourceFormat")
	targetFormat := r.FormValue("targetFormat")
	filename := strings.TrimSuffix(header.Filename, "."+sourceFormat)
	ratio, err := strconv.Atoi(r.FormValue("ratio"))
	if err != nil {
		reportError(w, fmt.Errorf("invalid ratio form value"), http.StatusBadRequest)
		return
	}

	if err = validation.ValidateConversionRequest(filename, sourceFormat, targetFormat, ratio); err != nil {
		reportError(w, err, http.StatusBadRequest)
		return
	}

	userID, ok := appcontext.UserIDFromContext(r.Context())
	if !ok {
		reportError(w, fmt.Errorf("can't get user id from application context"), http.StatusInternalServerError)
		return
	}

	sourceFileID, err := s.repo.InsertImage(filename, sourceFormat)
	if err != nil {
		reportError(w, fmt.Errorf("repository error: %w", err), http.StatusInternalServerError)
		return
	}
	logger.FromContext(r.Context()).WithField("file_id", sourceFileID).
		Infoln("original file successfully saved in the database")

	err = s.storage.UploadFile(sourceFile, sourceFileID)
	if err != nil {
		reportError(w, fmt.Errorf("storage error: %w", err), http.StatusInternalServerError)
		return
	}
	logger.FromContext(r.Context()).WithField("file_id", sourceFileID).
		Infoln("original file successfully uploaded to the S3 storage")

	requestID, err := s.repo.MakeRequest(userID, sourceFileID, sourceFormat, targetFormat, ratio)
	if err != nil {
		reportError(w, fmt.Errorf("repository error: %w", err), http.StatusInternalServerError)
		return
	}
	logger.FromContext(r.Context()).WithField("request_id", requestID).
		Infoln("request created with the status \"queued\"")

	sendResponse(w, ConvertResponse{RequestID: requestID}, http.StatusAccepted)

	err = s.producer.SendToQueue(sourceFileID, filename, sourceFormat, targetFormat, requestID, ratio)
	if err != nil {
		reportError(w, fmt.Errorf("can't send data to queue: %w", err), http.StatusInternalServerError)
		return
	}
	logger.FromContext(r.Context()).Infoln("message has been sent to the queue")
}

// DownloadImage allows you to download original/converted image by id.
func (s *Server) DownloadImage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(IDQueryKey)
	if id == "" {
		reportError(w, fmt.Errorf("id is missing in parameters"), http.StatusBadRequest)
		return
	}

	logger.FromContext(r.Context()).WithField("file_id", id).Infoln("get file_id from the URL")

	imageID, err := s.repo.GetImageIDInStore(id)
	if errors.Is(err, repository.ErrNoSuchImage) {
		reportError(w, err, http.StatusNotFound)
		return
	}
	if err != nil {
		reportError(w, err, http.StatusInternalServerError)
		return
	}

	url, err := s.storage.GetDownloadURL(imageID)
	if err != nil {
		reportError(w, err, http.StatusInternalServerError)
		return
	}

	sendResponse(w, DownloadResponse{ImageURL: url}, http.StatusOK)
}

// GetRequestsHistory displays the user's request history.
func (s *Server) GetRequestsHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := appcontext.UserIDFromContext(r.Context())
	if !ok {
		reportError(w, fmt.Errorf("can't get user id from application contex"), http.StatusInternalServerError)
		return
	}

	requestsHistory, err := s.repo.GetRequestsByUserID(userID)
	if err != nil {
		reportError(w, err, http.StatusInternalServerError)
		return
	}

	sendResponse(w, requestsHistory, http.StatusOK)
}
