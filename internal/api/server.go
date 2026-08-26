package api

import (
	"encoding/json"
	"net/http"
	"pharmacy-counter/internal/service"
	"strings"
)

type Server struct{ Pharmacy *service.Pharmacy }

func New(p *service.Pharmacy) *Server { return &Server{Pharmacy: p} }
func (s *Server) Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/api/board", s.board)
	m.HandleFunc("/api/tickets/call", s.call)
	m.HandleFunc("/api/tickets/", s.ticket)
	m.HandleFunc("/", s.index)
	return m
}
func write(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
func (s *Server) board(w http.ResponseWriter, r *http.Request) {
	b, e := s.Pharmacy.Board()
	if e != nil {
		write(w, map[string]string{"error": e.Error()}, 500)
		return
	}
	write(w, b, 200)
}
func (s *Server) call(w http.ResponseWriter, r *http.Request) {
	t, e := s.Pharmacy.CallNext()
	if e != nil {
		write(w, map[string]string{"error": e.Error()}, 400)
		return
	}
	write(w, t, 200)
}
func (s *Server) ticket(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tickets/")
	t, e := s.Pharmacy.Ticket(id)
	if e != nil {
		write(w, map[string]string{"error": e.Error()}, 404)
		return
	}
	write(w, t, 200)
}
func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte("<!doctype html><html><body><div id=app>Pharmacy Counter</div></body></html>"))
}
