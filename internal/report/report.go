package report

import (
	"encoding/json"
	"pharmacy-counter/internal/domain"
	"pharmacy-counter/internal/queue"
	"sort"
)

type Snapshot struct {
	Waiting   []domain.PickupTicket `json:"waiting"`
	Called    []domain.PickupTicket `json:"called"`
	Completed []domain.PickupTicket `json:"completed"`
	Total     int                   `json:"total"`
}

func FromBoard(b queue.Board) Snapshot {
	return Snapshot{Waiting: b.Waiting, Called: b.Called, Completed: b.Completed, Total: b.Total()}
}
func Encode(s Snapshot) ([]byte, error) { return json.MarshalIndent(s, "", "  ") }
func StatusLabel(status string) string {
	switch status {
	case domain.StatusWaiting:
		return "待取药"
	case domain.StatusCalled:
		return "已叫号"
	case domain.StatusCompleted:
		return "已完成"
	default:
		return "未知"
	}
}
func SortByNumber(ts []domain.PickupTicket) []domain.PickupTicket {
	out := append([]domain.PickupTicket(nil), ts...)
	sort.Slice(out, func(i, j int) bool { return out[i].Number < out[j].Number })
	return out
}
func CountByStatus(ts []domain.PickupTicket) map[string]int {
	m := map[string]int{}
	for _, t := range ts {
		m[t.Status]++
	}
	return m
}
