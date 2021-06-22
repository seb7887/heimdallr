package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/seb7887/heimdallr/health"
	log "github.com/sirupsen/logrus"
)

type HttpServer interface {
	Router() http.Handler
	GetHealth(w http.ResponseWriter, r *http.Request)
	Serve(ctx context.Context) error
}

type httpServer struct {
	router   http.Handler
	service  health.Service
	httpAddr string
}

func New(s health.Service, addr string) HttpServer {
	server := &httpServer{
		service:  s,
		httpAddr: addr,
	}
	router(server)

	return server
}

func router(s *httpServer) {
	r := mux.NewRouter()

	r.HandleFunc("/health", s.GetHealth).Methods(http.MethodGet)

	s.router = r
}

func (s *httpServer) Router() http.Handler {
	return s.router
}

func (s *httpServer) GetHealth(w http.ResponseWriter, r *http.Request) {
	err := s.service.GetHealth(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "ok")
}

func (s *httpServer) Serve(ctx context.Context) error {
	return http.ListenAndServe(s.httpAddr, s.Router())
}
