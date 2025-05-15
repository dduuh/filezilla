package handlers

import (
	"context"
	"filezilla/internal/domain"
	"filezilla/pkg/responses"
	"net/http"
	"strings"
)

func (h *Handler) userIdentify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			responses.HTTPError(w, "no such cookie header", http.StatusUnauthorized)
			return
		}

		token := strings.Split(header, " ")
		if len(token) != 2 {
			responses.HTTPError(w, "invalid token", http.StatusUnauthorized)
			return
		}

		if token[0] != "Bearer" {
			responses.HTTPError(w, "expected 'Bearer '", http.StatusUnauthorized)
			return
		}

		if token[1] == "" {
			responses.HTTPError(w, "invalid token", http.StatusUnauthorized)
			return
		}

		userId, err := h.services.Authentication.ParseToken(token[1])
		if err != nil {
			responses.HTTPError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), domain.UserIDKey, userId)
		r = r.WithContext(ctx)
		
		next.ServeHTTP(w, r)
	})
}
