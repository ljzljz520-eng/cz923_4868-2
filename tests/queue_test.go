package tests

import (
	"pharmacy-counter/internal/domain"
	"pharmacy-counter/internal/queue"
	"testing"
)

func TestQueueBoard(t *testing.T) {
	b := queue.BuildBoard([]domain.PickupTicket{{ID: "2", Number: 2, Status: domain.StatusWaiting}, {ID: "1", Number: 1, Status: domain.StatusWaiting}})
	if b.Waiting[0].Number != 1 {
		t.Fatal(b)
	}
}
