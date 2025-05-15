package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const ttl = 10 * time.Second

type Server struct {
	srv *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.srv = &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        handler,
		ReadTimeout:    ttl,
		WriteTimeout:   ttl,
		MaxHeaderBytes: 1 << 20,
	}

	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Shutdown(ctx)
}
