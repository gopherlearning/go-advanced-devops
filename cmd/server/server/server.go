package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gopherlearning/go-advanced-devops/cmd/server/handlers"
)

type Server struct {
	serv *http.Server
}

func NewServer(listen string, h *handlers.Handler) *Server {
	mux := http.NewServeMux()
	mux.Handle("/update/", http.HandlerFunc(h.Update))
	return &Server{serv: &http.Server{
		Addr:    listen,
		Handler: mux,
	}}
}
func (s *Server) Start() error {
	err := s.serv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return s.serv.Shutdown(ctx)
}
