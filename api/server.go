package api

import (
	"fmt"
	"net/http"
)

type Server struct {
	s *http.Server
}

func NewServer(port uint16, handler http.Handler) *Server {
	s := &Server{
		s: &http.Server{
			Handler: handler,
			Addr:    fmt.Sprintf(":%d", port),
		},
	}

	return s
}

func (s *Server) Listen() {
	err := s.s.ListenAndServe()

	if err != nil {
		// TODO: replace with proper logger
		fmt.Printf("Error occured while starting server, %s\n", err.Error())
	}
}
