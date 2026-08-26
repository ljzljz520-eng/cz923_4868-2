package tests

import (
	"pharmacy-counter/internal/domain"
	"testing"
)

func TestBusinessChain05(t *testing.T) {
	p := makePharmacy(t)
	_ = p.RegisterPatient(domain.Patient{ID: "p5", Name: "钱七"})
	_ = p.CreatePrescription(domain.Prescription{ID: "rx5", PatientID: "p5", Drug: "药", Quantity: 1})
	ticket, _ := p.CreatePickup("p5", "rx5")
	_, _ = p.CallNext()
	d, e := p.CompleteTicket(ticket.ID, "药师", 0)
	if e == nil {
		t.Errorf("expected validation failure, got success with %+v", d)
	}
}
