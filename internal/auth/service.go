package auth

import (
	"context"
	"errors"

	"company-service/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Repository defines interface of user repository
type Repository interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Login(ctx context.Context, username, password string) (string, string, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (string, error) {
	userID, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}
	return jwt.GenerateAccessToken(userID)
}
