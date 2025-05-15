package handlers

import (
	"filezilla/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	services service.Service
}

func NewHandler(services service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitHandler() *mux.Router {
	r := mux.NewRouter()

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/signup", h.signUp).Methods("POST")
	auth.HandleFunc("/login", h.login).Methods("POST")
	auth.HandleFunc("/refresh", h.refresh).Methods("POST")

	fileRouter := r.PathPrefix("/api/v1").Subrouter()
	fileRouter.Use(h.userIdentify)

	fileRouter.HandleFunc("/upload", h.uploadFile).Methods("POST")
	fileRouter.HandleFunc("/files", h.getFiles).Methods("GET")

	return r
}
