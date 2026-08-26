package domain

import "strings"

func NormalizeName(v string) string  { return strings.Join(strings.Fields(strings.TrimSpace(v)), " ") }
func NormalizePhone(v string) string { return strings.TrimSpace(v) }
func ActivePatients(values []Patient) []Patient {
	out := []Patient{}
	for _, v := range values {
		if v.Active {
			out = append(out, v)
		}
	}
	return out
}
