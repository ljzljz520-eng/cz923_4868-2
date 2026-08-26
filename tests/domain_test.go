package tests

import (
	"pharmacy-counter/internal/domain"
	"testing"
)

func TestDomainTransitions(t *testing.T) {
	v := domain.PickupTicket{ID: "t", PrescriptionID: "r", Number: 1, Status: domain.StatusWaiting}
	if e := v.Call("a"); e != nil {
		t.Fatal(e)
	}
	if e := v.Complete("b"); e != nil {
		t.Fatal(e)
	}
	if !v.IsCompleted() {
		t.Fatal("not completed")
	}
}
