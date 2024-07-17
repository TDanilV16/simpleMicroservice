package server

import (
	"net/http"
	"time"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return server.ListenAndServe()
}
