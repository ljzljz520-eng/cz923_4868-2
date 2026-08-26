package service

import (
	"errors"
	"pharmacy-counter/internal/domain"
)

func (p *Pharmacy) PreparePickup(patient domain.Patient, prescription domain.Prescription) (domain.PickupTicket, error) {
	if e := p.RegisterPatient(patient); e != nil {
		if _, found := p.FindPatient(patient.ID); found != nil {
			return domain.PickupTicket{}, e
		}
	}
	if e := p.CreatePrescription(prescription); e != nil {
		return domain.PickupTicket{}, e
	}
	if e := p.MarkReady(prescription.ID); e != nil {
		return domain.PickupTicket{}, e
	}
	return p.CreatePickup(patient.ID, prescription.ID)
}
func (p *Pharmacy) FinishCalledTicket(ticketID, pharmacist string, quantity int) (domain.DispenseRecord, error) {
	if e := EnsurePharmacist(pharmacist); e != nil {
		return domain.DispenseRecord{}, e
	}
	if e := EnsureQuantity(quantity); e != nil {
		return domain.DispenseRecord{}, e
	}
	return p.CompleteTicket(ticketID, pharmacist, quantity)
}
func (p *Pharmacy) RequireReadyPrescription(id string) error {
	v, e := p.Store.GetPrescription(id)
	if e != nil {
		return e
	}
	if !v.Ready {
		return errors.New("prescription not ready")
	}
	return nil
}
