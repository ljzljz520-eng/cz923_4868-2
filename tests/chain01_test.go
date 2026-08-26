package tests

import (
	"pharmacy-counter/internal/clock"
	"pharmacy-counter/internal/domain"
	"pharmacy-counter/internal/service"
	"pharmacy-counter/internal/storage"
	"testing"
)

func makePharmacy(t *testing.T) *service.Pharmacy {
	t.Helper()
	s, e := storage.Open(t.TempDir() + "/db")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { s.Close() })
	return service.New(s, clock.New("2026-01-01T00:00:00Z"))
}
func TestBusinessChain01(t *testing.T) {
	p := makePharmacy(t)
	if e := p.RegisterPatient(domain.Patient{ID: "p1", Name: "张三", Active: true}); e != nil {
		t.Fatal(e)
	}
	if e := p.CreatePrescription(domain.Prescription{ID: "rx1", PatientID: "p1", Drug: "阿莫西林", Dosage: "500mg", Quantity: 2}); e != nil {
		t.Fatal(e)
	}
	if e := p.MarkReady("rx1"); e != nil {
		t.Fatal(e)
	}
}
