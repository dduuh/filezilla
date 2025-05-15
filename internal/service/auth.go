package service

import (
	"errors"
	"filezilla/internal/domain"
	"filezilla/internal/repository"
	"filezilla/pkg/hash"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type AuthService struct {
	repo repository.Authentication
}

func NewAuthService(repo repository.Authentication) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) Create(user domain.User) (int, error) {
	user.Password = hash.HashPassword(user.Password)
	return a.repo.Create(user)
}

func (a *AuthService) NewJWT(userId int) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(viper.GetDuration("auth.accessTtl")).Unix(),
		Subject:   strconv.Itoa(userId),
	})

	return accessToken.SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
}

func (a *AuthService) CreateToken(email, password string) (string, int, error) {
	hashedPassword := hash.HashPassword(password)

	userId, err := a.repo.GetUser(email, hashedPassword)
	if err != nil {
		return "", -1, err
	}

	token, err := a.NewJWT(userId)
	if err != nil {
		return "", -1, err
	}

	return token, userId, nil
}

func (a *AuthService) NewRefreshToken(userId int) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(viper.GetDuration("auth.refreshTtl")).Unix(),
		Subject:   strconv.Itoa(userId),
	})

	return accessToken.SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
}

func (a *AuthService) CreateRefreshToken(userId int) (string, error) {
	refreshToken, err := a.NewRefreshToken(userId)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (a *AuthService) ParseToken(token string) (int, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
	})
	if err != nil {
		return -1, err
	}

	claims, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return -1, errors.New("claims not found")
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (a *AuthService) RefreshTokens(cookie *http.Cookie, token string) (string, string, error) {
	userId, err := a.ParseToken(token)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err := a.NewJWT(userId)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := a.NewRefreshToken(userId)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
