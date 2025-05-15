package service

import (
	"context"
	"filezilla/internal/domain"
	"filezilla/internal/repository"
	"net/http"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authentication interface {
	Create(user domain.User) (int, error)

	NewJWT(userId int) (string, error)
	NewRefreshToken(userId int) (string, error)

	ParseToken(token string) (int, error)

	CreateRefreshToken(userId int) (string, error)           // Uses: NewRefreshToken()
	CreateToken(email, password string) (string, int, error) // Uses NewJWT()

	RefreshTokens(cookie *http.Cookie, token string) (string, string, error)
}

type Files interface {
	CreateFile(file domain.File) (int, error)
	GetFiles() ([]domain.File, error)
}

type Log interface {
	Log(ctx context.Context, logInput *domain.Log) error
}

type Service struct {
	Authentication
	Files
	Log
}

func NewService(repos repository.Repository) *Service {
	return &Service{
		Authentication: NewAuthService(repos.Authentication),
		Files:          NewFileService(repos.Files),
		Log:            NewLogService(repos.Logs),
	}
}
