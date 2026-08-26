package service

import (
	"pharmacy-counter/internal/domain"
	"strings"
)

func (p *Pharmacy) SearchPatients(term string) ([]domain.Patient, error) {
	all, e := p.Patients()
	if e != nil {
		return nil, e
	}
	term = strings.ToLower(strings.TrimSpace(term))
	out := []domain.Patient{}
	for _, v := range all {
		if term == "" || strings.Contains(strings.ToLower(v.Name), term) || strings.Contains(v.Phone, term) {
			out = append(out, v)
		}
	}
	return out, nil
}
func (p *Pharmacy) TicketsForPatient(patientID string) ([]domain.PickupTicket, error) {
	all, e := p.Tickets()
	if e != nil {
		return nil, e
	}
	ps, e := p.Prescriptions()
	if e != nil {
		return nil, e
	}
	ids := map[string]bool{}
	for _, v := range ps {
		if v.PatientID == patientID {
			ids[v.ID] = true
		}
	}
	out := []domain.PickupTicket{}
	for _, v := range all {
		if ids[v.PrescriptionID] {
			out = append(out, v)
		}
	}
	return out, nil
}
func (p *Pharmacy) PendingCount() (int, error) {
	b, e := p.Board()
	if e != nil {
		return 0, e
	}
	return len(b.Waiting), nil
}
