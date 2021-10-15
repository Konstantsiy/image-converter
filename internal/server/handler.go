// Package server implements http handlers.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Konstantsiy/image-converter/internal/queue"

	"github.com/Konstantsiy/image-converter/internal/service"

	"github.com/Konstantsiy/image-converter/pkg/logger"

	"github.com/Konstantsiy/image-converter/internal/validation"
	"github.com/gorilla/mux"
)

const IDQueryKey = "id"

// Server represents application server.
type Server struct {
	authService     service.Authorization
	imageService    service.Images
	requestsService service.Requests
	producer        *queue.Producer
}

func NewServer(authService service.Authorization, imageService service.Images, requestsService service.Requests, producer *queue.Producer) *Server {
	return &Server{
		authService:     authService,
		imageService:    imageService,
		requestsService: requestsService,
		producer:        producer}
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
	type authRequest struct {
		Email    string
		Password string
	}

	var request authRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		reportErrorWithCode(w, err, http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()

	if request.Email == "" || request.Password == "" {
		reportErrorWithCode(w, fmt.Errorf("empty field"), http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := s.authService.LogIn(r.Context(), request.Email, request.Password)
	if err != nil {
		reportError(w, err)
	}

	type loginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	sendResponse(w, loginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, http.StatusOK)
}

// SignUp implements the user registration process.
func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	type authRequest struct {
		Email    string
		Password string
	}

	var request authRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		reportErrorWithCode(w, err, http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()

	if err = validation.ValidateSignUpRequest(request.Email, request.Password); err != nil {
		reportErrorWithCode(w, err, http.StatusBadRequest)
		return
	}

	userID, err := s.authService.SignUp(r.Context(), request.Email, request.Password)
	if err != nil {
		reportError(w, err)
		return
	}

	type signUpResponse struct {
		UserID string `json:"user_id"`
	}

	sendResponse(w, signUpResponse{UserID: userID}, http.StatusOK)
}

// ConvertImage converts needed image according to the request.
func (s *Server) ConvertImage(w http.ResponseWriter, r *http.Request) {
	sourceFile, header, err := r.FormFile("file")
	if err != nil {
		reportErrorWithCode(w, fmt.Errorf("can't get sourceFile from form"), http.StatusBadRequest)
		return
	}
	defer sourceFile.Close()

	sourceFormat := r.FormValue("sourceFormat")
	targetFormat := r.FormValue("targetFormat")
	filename := strings.TrimSuffix(header.Filename, "."+sourceFormat)
	ratio, err := strconv.Atoi(r.FormValue("ratio"))
	if err != nil {
		reportErrorWithCode(w, fmt.Errorf("invalid ratio form value"), http.StatusBadRequest)
		return
	}

	if err = validation.ValidateConversionRequest(filename, sourceFormat, targetFormat, ratio); err != nil {
		reportErrorWithCode(w, err, http.StatusBadRequest)
		return
	}

	sourceFileID, requestID, err := s.imageService.Convert(r.Context(), sourceFile, filename, sourceFormat, targetFormat, ratio)
	if err != nil {
		reportError(w, err)
		return
	}

	type convertResponse struct {
		RequestID string `json:"request_id"`
	}

	sendResponse(w, convertResponse{RequestID: requestID}, http.StatusAccepted)

	err = s.producer.SendToQueue(sourceFileID, filename, sourceFormat, targetFormat, requestID, ratio)
	if err != nil {
		reportErrorWithCode(w, fmt.Errorf("can't send data to queue: %w", err), http.StatusInternalServerError)
		return
	}
	logger.FromContext(r.Context()).Infoln("message has been sent to the queue")
}

// DownloadImage allows you to download original/converted image by id.
func (s *Server) DownloadImage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(IDQueryKey)
	if id == "" {
		reportErrorWithCode(w, fmt.Errorf("image id is missing in parameters"), http.StatusBadRequest)
		return
	}

	url, err := s.imageService.Download(r.Context(), id)
	if err != nil {
		reportError(w, err)
		return
	}

	type downloadResponse struct {
		ImageURL string `json:"image_url"`
	}

	sendResponse(w, downloadResponse{ImageURL: url}, http.StatusOK)
}

// GetRequestsHistory displays the user's request history.
func (s *Server) GetRequestsHistory(w http.ResponseWriter, r *http.Request) {
	requests, err := s.requestsService.GetUsersRequests(r.Context())
	if err != nil {
		reportError(w, err)
		return
	}

	sendResponse(w, requests, http.StatusOK)
}
