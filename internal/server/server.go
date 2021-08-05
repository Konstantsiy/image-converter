// Package server implements http handlers.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Konstantsiy/image-converter/internal/converter"
	"github.com/Konstantsiy/image-converter/internal/appcontext"
	"github.com/Konstantsiy/image-converter/internal/auth"
	"github.com/Konstantsiy/image-converter/internal/hash"
	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/internal/validation"
	"github.com/gorilla/mux"
)

const (
	AuthorizationHeader  = "Authorization"
	NeededSecurityScheme = "Bearer"
)

// AuthRequest represents the user's authorization request.
type AuthRequest struct {
	Email    string
	Password string
}

// ConversionRequest represents an image conversion request.
type ConversionRequest struct {
	File         string
	SourceFormat string
	TargetFormat string
	Ratio        int
}

// LoginResponse represents token for authorization response.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

//SignUpResponse represents user id from sign up response.
type SignUpResponse struct {
	UserID string `json:"user_id"`
}

//DownloadResponse represents downloaded image URL.
type DownloadResponse struct {
	ImageURL string `json:"image_url"`
}

// Server represents application server.
type Server struct {
	repo         *repository.Repository
	tokenManager *auth.TokenManager
	conv         *converter.Converter
}

// NewServer creates new application server.
func NewServer(repo *repository.Repository, tokenManager *auth.TokenManager, conv *converter.Converter) *Server {
	return &Server{
		repo:         repo,
		tokenManager: tokenManager,
		conv:         conv,
	}
}

// AuthMiddleware checks user authorization.
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(AuthorizationHeader)
		if authHeader == "" {
			http.Error(w, "empty auth handler", http.StatusUnauthorized)
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != NeededSecurityScheme {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		if len(authHeaderParts[1]) == 0 {
			http.Error(w, "token is empty", http.StatusUnauthorized)
			return
		}

		token := authHeaderParts[1]
		claimsUserID, err := s.tokenManager.ParseToken(token)
		if err != nil {
			http.Error(w, "can't parse JWT", http.StatusUnauthorized)
			return
		}

		ctx := appcontext.ContextWithUserID(r.Context(), claimsUserID)

		next.ServeHTTP(w, r.WithContext(ctx))
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

	fmt.Fprint(w, &LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken})
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
	if err == repository.ErrUserAlreadyExists {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, &SignUpResponse{UserID: userID})
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

	convFile, err := s.conv.Convert(file, targetFormat, ratio)
	if err != nil {
		http.Error(w, "can't convert image:"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(convFile.Name())

	// ... next PRs
}

// DownloadImage allows you to download original/converted image by id.
func (s *Server) DownloadImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := vars["id"]

	image, err := s.repo.GetImageByID(imageID)
	if err == repository.ErrNoSuchImage {
		http.Error(w, "can't get image info: "+err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get image downloaded URL from storage by image.Location
	url := "http(s)://s3.amazonaws.com/" + image.Location + "/file_name.extension"

	fmt.Fprint(w, &DownloadResponse{ImageURL: url})
}

// GetRequestsHistory displays the user's request history.
func (s *Server) GetRequestsHistory(w http.ResponseWriter, r *http.Request) {
	userID := appcontext.UserIDFromContext(r.Context())

	requestsHistory, err := s.repo.GetRequestsByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, fmt.Sprint(requestsHistory))
}
