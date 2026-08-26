package tests

import (
	"pharmacy-counter/internal/domain"
	"testing"
)

func TestBusinessChain02(t *testing.T) {
	p := makePharmacy(t)
	_ = p.RegisterPatient(domain.Patient{ID: "p2", Name: "李四"})
	_ = p.CreatePrescription(domain.Prescription{ID: "rx2", PatientID: "p2", Drug: "维生素", Quantity: 1})
	ticket, e := p.CreatePickup("p2", "rx2")
	if e != nil {
		t.Fatal(e)
	}
	if ticket.Number != 1 || !ticket.IsWaiting() {
		t.Fatalf("unexpected ticket %+v", ticket)
	}
}
