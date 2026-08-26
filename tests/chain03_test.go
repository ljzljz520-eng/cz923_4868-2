package tests

import (
	"pharmacy-counter/internal/domain"
	"testing"
)

func TestBusinessChain03(t *testing.T) {
	p := makePharmacy(t)
	_ = p.RegisterPatient(domain.Patient{ID: "p3", Name: "王五"})
	_ = p.CreatePrescription(domain.Prescription{ID: "rx3", PatientID: "p3", Drug: "药", Quantity: 1})
	ticket, _ := p.CreatePickup("p3", "rx3")
	called, e := p.CallNext()
	if e != nil || called.ID != ticket.ID {
		t.Fatal(e)
	}
	d, e := p.FinishCalledTicket(ticket.ID, "药师", 1)
	if e != nil {
		t.Fatal(e)
	}
	if d.Quantity != 1 {
		t.Fatal(d)
	}
}
