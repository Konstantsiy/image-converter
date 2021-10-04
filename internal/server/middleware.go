package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Konstantsiy/image-converter/internal/appcontext"
	"github.com/Konstantsiy/image-converter/pkg/logger"
)

// StatusRecorder contains a writer for storing the requests status code.
type StatusRecorder struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader saves requests status code.
func (sr *StatusRecorder) WriteHeader(statusCode int) {
	sr.StatusCode = statusCode
	sr.ResponseWriter.WriteHeader(statusCode)
}

// LoggingMiddleware logs http requests after they are executed.
func (s *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := logger.ContextWithLogger(r.Context())

		recorder := &StatusRecorder{
			ResponseWriter: w,
			StatusCode:     DefaultStatusCode,
		}

		next.ServeHTTP(recorder, r.WithContext(ctx))

		logger.CompleteRequest(ctx, r, time.Since(start), recorder.StatusCode)
	})
}

// AuthMiddleware checks user authorization.
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(AuthorizationHeader)
		if authHeader == "" {
			reportError(w, fmt.Errorf("empty auth handler"), http.StatusUnauthorized)
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != NeededSecurityScheme {
			reportError(w, fmt.Errorf("invalid auth handler"), http.StatusUnauthorized)
			return
		}

		if len(authHeaderParts[1]) == 0 {
			reportError(w, fmt.Errorf("token is empty"), http.StatusUnauthorized)
			return
		}

		token := authHeaderParts[1]
		claimsUserID, err := s.tokenManager.ParseToken(token)
		if err != nil {
			reportError(w, fmt.Errorf("can't parse JWT: %w", err), http.StatusUnauthorized)
			return
		}

		ctx := appcontext.ContextWithUserID(r.Context(), claimsUserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
