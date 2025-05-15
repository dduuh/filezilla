package handlers

import (
	"context"
	"encoding/json"
	"filezilla/internal/domain"
	cook "filezilla/pkg/cookie"
	"filezilla/pkg/responses"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var user domain.User

	if r.Method != http.MethodPost {
		responses.HTTPError(w, "incorrect HTTP method", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		responses.HTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := domain.Validate.Struct(user); err != nil {
		responses.HTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := h.services.Authentication.Create(user)
	if err != nil {
		responses.HTTPError(w, "failed to create a user", http.StatusInternalServerError)
		return
	}

	if err := h.services.Logs.Log(context.TODO(), &domain.Log{
		UserId: userId,
		Role:   domain.USER,
		Action: domain.REGISTER,
	}); err != nil {
		logrus.Fatal(err.Error())
	}

	responses.HTTPResponse(w, http.StatusCreated, "id", userId)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var user domain.User

	if r.Method != http.MethodPost {
		responses.HTTPError(w, "incorrect HTTP method", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		responses.HTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := domain.Validate.Struct(user); err != nil {
		responses.HTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create JWT Access and Refresh token
	accessToken, userId, err := h.services.Authentication.CreateToken(user.Email, user.Password)
	if err != nil {
		responses.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.services.Authentication.CreateRefreshToken(userId)
	if err != nil {
		responses.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("refresh-token", refreshToken)

	cook.SetCookie(w, "refresh-token", refreshToken, "/", "", true, true, http.SameSiteLaxMode)

	if err := h.services.Logs.Log(context.TODO(), &domain.Log{
		UserId: userId,
		Role:   domain.USER,
		Action: domain.LOGIN,
	}); err != nil {
		logrus.Fatal(err.Error())
	}

	responses.HTTPResponse(w, http.StatusOK, "token", accessToken)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.HTTPError(w, "incorrect HTTP method", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		responses.HTTPError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	refreshToken := cookie.Value

	newAccessToken, newRefreshToken, err := h.services.Authentication.RefreshTokens(cookie, refreshToken)
	if err != nil {
		responses.HTTPError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	cook.SetCookie(w, "refresh-token", newRefreshToken, "/", "", true, true, http.SameSiteLaxMode)

	responses.HTTPResponse(w, http.StatusOK, "New token", newAccessToken)
}
