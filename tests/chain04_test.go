package tests

import (
	"pharmacy-counter/internal/domain"
	"pharmacy-counter/internal/report"
	"testing"
)

func TestBusinessChain04(t *testing.T) {
	p := makePharmacy(t)
	_ = p.RegisterPatient(domain.Patient{ID: "p4", Name: "赵六"})
	_ = p.CreatePrescription(domain.Prescription{ID: "rx4", PatientID: "p4", Drug: "药", Quantity: 1})
	_, _ = p.CreatePickup("p4", "rx4")
	b, e := p.Board()
	if e != nil {
		t.Fatal(e)
	}
	if report.Headline(report.FromBoard(b)) == "" {
		t.Fatal("empty headline")
	}
}
