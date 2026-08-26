package tests

import (
	"pharmacy-counter/internal/clock"
	"pharmacy-counter/internal/domain"
	"pharmacy-counter/internal/service"
	"pharmacy-counter/internal/storage"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/db"
	s, e := storage.Open(path)
	if e != nil {
		t.Fatal(e)
	}
	p := service.New(s, clock.New("t"))
	_ = p.RegisterPatient(domain.Patient{ID: "persist", Name: "持久"})
	_ = s.Close()
	s, e = storage.Open(path)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	if _, e = s.GetPatient("persist"); e != nil {
		t.Fatal(e)
	}
}
