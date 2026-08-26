package queue

import "pharmacy-counter/internal/domain"

func WaitingOnly(values []domain.PickupTicket) []domain.PickupTicket {
	out := []domain.PickupTicket{}
	for _, v := range values {
		if v.Status == domain.StatusWaiting {
			out = append(out, v)
		}
	}
	return out
}
func CalledOnly(values []domain.PickupTicket) []domain.PickupTicket {
	out := []domain.PickupTicket{}
	for _, v := range values {
		if v.Status == domain.StatusCalled {
			out = append(out, v)
		}
	}
	return out
}
func CompletedOnly(values []domain.PickupTicket) []domain.PickupTicket {
	out := []domain.PickupTicket{}
	for _, v := range values {
		if v.Status == domain.StatusCompleted {
			out = append(out, v)
		}
	}
	return out
}
