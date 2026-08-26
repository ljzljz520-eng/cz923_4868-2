package queue

import (
	"pharmacy-counter/internal/domain"
	"sort"
)

type Board struct{ Waiting, Called, Completed []domain.PickupTicket }

func BuildBoard(tickets []domain.PickupTicket) Board {
	b := Board{}
	for _, t := range tickets {
		switch t.Status {
		case domain.StatusWaiting:
			b.Waiting = append(b.Waiting, t)
		case domain.StatusCalled:
			b.Called = append(b.Called, t)
		case domain.StatusCompleted:
			b.Completed = append(b.Completed, t)
		}
	}
	sort.Slice(b.Waiting, func(i, j int) bool { return b.Waiting[i].Number < b.Waiting[j].Number })
	sort.Slice(b.Called, func(i, j int) bool { return b.Called[i].Number < b.Called[j].Number })
	sort.Slice(b.Completed, func(i, j int) bool { return b.Completed[i].Number > b.Completed[j].Number })
	return b
}
func (b Board) Total() int         { return len(b.Waiting) + len(b.Called) + len(b.Completed) }
func (b Board) HasWaiting() bool   { return len(b.Waiting) > 0 }
func (b Board) HasCalled() bool    { return len(b.Called) > 0 }
func (b Board) HasCompleted() bool { return len(b.Completed) > 0 }
