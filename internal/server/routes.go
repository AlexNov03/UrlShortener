package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) InitRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/shorten", s.delivery.ShortenUrl).Methods(http.MethodPost)
	router.HandleFunc("/{shortened_url}", s.delivery.GetOriginalUrl).Methods(http.MethodGet)
	s.server.Handler = router
}
