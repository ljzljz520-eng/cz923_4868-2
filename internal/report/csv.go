package report

import (
	"fmt"
	"pharmacy-counter/internal/domain"
	"strings"
)

func CSV(values []domain.PickupTicket) string {
	rows := []string{"id,number,status,created_at"}
	for _, v := range values {
		rows = append(rows, fmt.Sprintf("%s,%d,%s,%s", v.ID, v.Number, v.Status, v.CreatedAt))
	}
	return strings.Join(rows, "\n")
}
func Headline(b Snapshot) string {
	return fmt.Sprintf("待取药 %d | 已叫号 %d | 已完成 %d", len(b.Waiting), len(b.Called), len(b.Completed))
}
