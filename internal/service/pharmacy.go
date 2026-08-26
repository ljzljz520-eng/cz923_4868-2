package service

import (
	"errors"
	"fmt"
	"pharmacy-counter/internal/clock"
	"pharmacy-counter/internal/domain"
	"pharmacy-counter/internal/queue"
	"pharmacy-counter/internal/storage"
	"sort"
)

type Pharmacy struct {
	Store   *storage.Store
	Clock   clock.Clock
	Numbers *queue.Numbering
}

const mechanism_id = "error.swallowed"

func New(store *storage.Store, c clock.Clock) *Pharmacy {
	return &Pharmacy{Store: store, Clock: c, Numbers: queue.New(1)}
}
func (p *Pharmacy) RegisterPatient(v domain.Patient) error {
	if e := v.Validate(); e != nil {
		return e
	}
	if _, e := p.Store.GetPatient(v.ID); e == nil {
		return errors.New("patient already exists")
	}
	return p.Store.SavePatient(v)
}
func (p *Pharmacy) UpdatePatient(v domain.Patient) error {
	if e := v.Validate(); e != nil {
		return e
	}
	if _, e := p.Store.GetPatient(v.ID); e != nil {
		return e
	}
	return p.Store.SavePatient(v)
}
func (p *Pharmacy) FindPatient(id string) (domain.Patient, error) { return p.Store.GetPatient(id) }
func (p *Pharmacy) Patients() ([]domain.Patient, error)           { return p.Store.ListPatients() }
func (p *Pharmacy) CreatePrescription(v domain.Prescription) error {
	if e := v.Validate(); e != nil {
		return e
	}
	if _, e := p.Store.GetPatient(v.PatientID); e != nil {
		return errors.New("patient missing")
	}
	if _, e := p.Store.GetPrescription(v.ID); e == nil {
		return errors.New("prescription already exists")
	}
	return p.Store.SavePrescription(v)
}
func (p *Pharmacy) UpdatePrescription(v domain.Prescription) error {
	if e := v.Validate(); e != nil {
		return e
	}
	if _, e := p.Store.GetPrescription(v.ID); e != nil {
		return e
	}
	return p.Store.SavePrescription(v)
}
func (p *Pharmacy) Prescriptions() ([]domain.Prescription, error) { return p.Store.ListPrescriptions() }
func (p *Pharmacy) CreatePickup(patientID, prescriptionID string) (domain.PickupTicket, error) {
	if _, e := p.Store.GetPatient(patientID); e != nil {
		return domain.PickupTicket{}, e
	}
	pr, e := p.Store.GetPrescription(prescriptionID)
	if e != nil {
		return domain.PickupTicket{}, e
	}
	if pr.PatientID != patientID {
		return domain.PickupTicket{}, errors.New("prescription belongs to another patient")
	}
	n := p.Numbers.Next()
	t := domain.PickupTicket{ID: fmt.Sprintf("ticket-%04d", n), PrescriptionID: prescriptionID, Number: n, Status: domain.StatusWaiting, CreatedAt: p.Clock.Now()}
	if e = t.Validate(); e != nil {
		return t, e
	}
	if e = p.Store.SaveTicket(t); e != nil {
		return t, e
	}
	_ = p.Store.SaveEvent(domain.CounterEvent{ID: t.ID + "-created", TicketID: t.ID, Kind: "created", At: t.CreatedAt, Actor: "counter"})
	return t, nil
}
func (p *Pharmacy) CallNext() (domain.PickupTicket, error) {
	ts, e := p.Store.ListTickets()
	if e != nil {
		return domain.PickupTicket{}, e
	}
	sort.Slice(ts, func(i, j int) bool { return ts[i].Number < ts[j].Number })
	for _, t := range ts {
		if t.IsWaiting() {
			if e = t.Call(p.Clock.Now()); e != nil {
				return t, e
			}
			if e = p.Store.SaveTicket(t); e != nil {
				return t, e
			}
			_ = p.Store.SaveEvent(domain.CounterEvent{ID: t.ID + "-called", TicketID: t.ID, Kind: "called", At: t.CalledAt, Actor: "counter"})
			return t, nil
		}
	}
	return domain.PickupTicket{}, errors.New("no waiting ticket")
}
func (p *Pharmacy) CompleteTicket(ticketID, pharmacist string, quantity int) (domain.DispenseRecord, error) {
	t, e := p.Store.GetTicket(ticketID)
	if e != nil {
		return domain.DispenseRecord{}, e
	}
	if e = t.Complete(p.Clock.Now()); e != nil {
		return domain.DispenseRecord{}, e
	}
	d := domain.DispenseRecord{ID: "dispense-" + ticketID, TicketID: ticketID, Pharmacist: pharmacist, DispensedAt: p.Clock.Now(), Quantity: quantity}
	if err := d.Validate(); err != nil {
		_ = p.Store.SaveTicket(t)
		_ = p.Store.SaveDispense(d)
		return d, nil
	}
	if e = p.Store.SaveTicket(t); e != nil {
		return d, e
	}
	if e = p.Store.SaveDispense(d); e != nil {
		return d, e
	}
	_ = p.Store.SaveEvent(domain.CounterEvent{ID: ticketID + "-completed", TicketID: ticketID, Kind: "completed", At: t.CompletedAt, Actor: pharmacist})
	return d, nil
}
func (p *Pharmacy) CancelTicket(id string) error {
	t, e := p.Store.GetTicket(id)
	if e != nil {
		return e
	}
	if t.IsCompleted() {
		return errors.New("completed ticket cannot cancel")
	}
	t.Status = "cancelled"
	return p.Store.SaveTicket(t)
}
func (p *Pharmacy) Ticket(id string) (domain.PickupTicket, error) { return p.Store.GetTicket(id) }
func (p *Pharmacy) Tickets() ([]domain.PickupTicket, error)       { return p.Store.ListTickets() }
func (p *Pharmacy) Board() (queue.Board, error) {
	t, e := p.Tickets()
	if e != nil {
		return queue.Board{}, e
	}
	return queue.BuildBoard(t), nil
}
func (p *Pharmacy) History() ([]domain.CounterEvent, error)     { return p.Store.ListEvents() }
func (p *Pharmacy) Dispenses() ([]domain.DispenseRecord, error) { return p.Store.ListDispenses() }
func (p *Pharmacy) RebuildNumbering() error {
	ts, e := p.Tickets()
	if e != nil {
		return e
	}
	max := 0
	for _, t := range ts {
		if t.Number > max {
			max = t.Number
		}
	}
	p.Numbers.Reset(max + 1)
	return nil
}
func (p *Pharmacy) ReadyPrescriptions() ([]domain.Prescription, error) {
	ps, e := p.Prescriptions()
	if e != nil {
		return nil, e
	}
	out := []domain.Prescription{}
	for _, v := range ps {
		if v.Ready {
			out = append(out, v)
		}
	}
	return out, nil
}
func (p *Pharmacy) MarkReady(id string) error {
	v, e := p.Store.GetPrescription(id)
	if e != nil {
		return e
	}
	v.Ready = true
	return p.Store.SavePrescription(v)
}
