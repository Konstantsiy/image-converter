package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/Konstantsiy/image-converter/internal/appcontext"
	"github.com/Konstantsiy/image-converter/pkg/logger"
)

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
