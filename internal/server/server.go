package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AlexNov03/UrlShortener/internal/bootstrap"
	"github.com/AlexNov03/UrlShortener/internal/delivery"
)

type Server struct {
	server   *http.Server
	cfg      *bootstrap.Config
	handler  http.Handler
	delivery *delivery.UrlDelivery
}

func NewServer(cfg *bootstrap.Config, delivery *delivery.UrlDelivery) *Server {
	return &Server{cfg: cfg, delivery: delivery}
}

func (s *Server) Init() {
	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", "0.0.0.0", s.cfg.Server.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.InitRoutes()
}

func (s *Server) Run() error {
	log.Printf("starting server, listening on addr %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("error while starting server: %v ", err)
	}
	return nil
}

func (s *Server) Stop() error {
	if err := s.server.Close(); err != nil {
		return fmt.Errorf("error while stopping server: %v ", err)
	}
	log.Printf("server stopped successfully")
	return nil
}
