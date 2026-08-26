package domain

import (
	"errors"
	"fmt"
	"strings"
)

type Patient struct {
	ID, Name, Phone string
	Active          bool
}
type Prescription struct {
	ID, PatientID, Drug, Dosage string
	Quantity                    int
	Ready                       bool
}
type PickupTicket struct {
	ID, PrescriptionID string
	Number             int
	Status             string
	CreatedAt          string
	CalledAt           string
	CompletedAt        string
}
type DispenseRecord struct {
	ID, TicketID, Pharmacist, DispensedAt string
	Quantity                              int
	Notes                                 string
}
type CounterEvent struct {
	ID, TicketID, Kind, At, Actor string
	Detail                        string
}

const (
	StatusWaiting   = "waiting"
	StatusCalled    = "called"
	StatusCompleted = "completed"
)

func (p Patient) Validate() error {
	if strings.TrimSpace(p.ID) == "" {
		return errors.New("patient id required")
	}
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("patient name required")
	}
	return nil
}
func (p Prescription) Validate() error {
	if p.ID == "" || p.PatientID == "" {
		return errors.New("prescription identity required")
	}
	if p.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if strings.TrimSpace(p.Drug) == "" {
		return errors.New("drug required")
	}
	return nil
}
func (t PickupTicket) Validate() error {
	if t.ID == "" || t.PrescriptionID == "" {
		return errors.New("ticket identity required")
	}
	if t.Number <= 0 {
		return errors.New("ticket number required")
	}
	if t.Status != StatusWaiting && t.Status != StatusCalled && t.Status != StatusCompleted {
		return fmt.Errorf("invalid status %s", t.Status)
	}
	return nil
}
func (r DispenseRecord) Validate() error {
	if r.ID == "" || r.TicketID == "" {
		return errors.New("dispense identity required")
	}
	if r.Quantity <= 0 {
		return errors.New("dispense quantity must be positive")
	}
	return nil
}
func (e CounterEvent) Validate() error {
	if e.ID == "" || e.TicketID == "" || e.Kind == "" {
		return errors.New("event identity required")
	}
	return nil
}
func (t PickupTicket) IsWaiting() bool   { return t.Status == StatusWaiting }
func (t PickupTicket) IsCalled() bool    { return t.Status == StatusCalled }
func (t PickupTicket) IsCompleted() bool { return t.Status == StatusCompleted }
func (t *PickupTicket) Call(at string) error {
	if !t.IsWaiting() {
		return errors.New("ticket is not waiting")
	}
	t.Status = StatusCalled
	t.CalledAt = at
	return nil
}
func (t *PickupTicket) Complete(at string) error {
	if !t.IsCalled() {
		return errors.New("ticket must be called first")
	}
	t.Status = StatusCompleted
	t.CompletedAt = at
	return nil
}
func (p Prescription) Summary() string       { return fmt.Sprintf("%s %s x%d", p.Drug, p.Dosage, p.Quantity) }
func (t PickupTicket) DisplayNumber() string { return fmt.Sprintf("%03d", t.Number) }
