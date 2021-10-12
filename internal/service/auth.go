package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/pkg/hash"
	"github.com/Konstantsiy/image-converter/pkg/jwt"
	"github.com/Konstantsiy/image-converter/pkg/logger"
)

type AuthService struct {
	usersRepo *repository.UsersRepository
	tm        *jwt.TokenManager
}

func NewAuthService(repo *repository.UsersRepository, tm *jwt.TokenManager) *AuthService {
	return &AuthService{usersRepo: repo, tm: tm}
}

func (auth *AuthService) LogIn(ctx context.Context, email, password string) (string, string, error) {
	user, err := auth.usersRepo.GetUserByEmail(ctx, email)
	if err == repository.ErrNoSuchUser {
		return "", "", &ServiceError{
			fmt.Errorf("invalid email or password"),
			http.StatusUnauthorized,
		}
	}
	if err != nil {
		return "", "", &ServiceError{err, http.StatusInternalServerError}
	}

	if ok, err := hash.ComparePasswordHash(password, user.Password); !ok || err != nil {
		return "", "", &ServiceError{
			fmt.Errorf("invalid email or password"),
			http.StatusUnauthorized,
		}
	}

	logger.FromContext(ctx).WithField("user_id", user.ID).Infoln("user successfully logged in")

	accessToken, err := auth.tm.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", &ServiceError{
			fmt.Errorf("can't generate access token: %w", err),
			http.StatusInternalServerError,
		}
	}

	refreshToken, err := auth.tm.GenerateRefreshToken()
	if err != nil {
		return "", "", &ServiceError{
			fmt.Errorf("can't generate refresh token: %w", err),
			http.StatusInternalServerError,
		}
	}

	return accessToken, refreshToken, nil
}

func (auth *AuthService) SignUp(ctx context.Context, email, password string) (string, error) {
	hashPwd, err := hash.GeneratePasswordHash(password)
	if err != nil {
		return "", &ServiceError{
			fmt.Errorf("can't generate password hash: %w", err),
			http.StatusInternalServerError,
		}
	}

	userID, err := auth.usersRepo.InsertUser(ctx, email, hashPwd)
	if errors.Is(err, repository.ErrUserAlreadyExists) {
		return "", &ServiceError{err, http.StatusBadRequest}
	}
	if err != nil {
		return "", &ServiceError{err, http.StatusInternalServerError}
	}

	return userID, nil
}

func (auth *AuthService) ParseToken(accessToken string) (string, error) {
	return auth.tm.ParseToken(accessToken)
}