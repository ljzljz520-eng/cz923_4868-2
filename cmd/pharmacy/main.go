package main

import (
	"log"
	"net/http"
	"os"
	"pharmacy-counter/internal/api"
	"pharmacy-counter/internal/clock"
	"pharmacy-counter/internal/service"
	"pharmacy-counter/internal/storage"
)

func main() {
	path := os.Getenv("PHARMACY_DB")
	if path == "" {
		path = "pharmacy.db"
	}
	s, e := storage.Open(path)
	if e != nil {
		log.Fatal(e)
	}
	defer s.Close()
	p := service.New(s, clock.New("2026-01-01T00:00:00Z"))
	log.Println("pharmacy counter listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", api.New(p).Handler()))
}
