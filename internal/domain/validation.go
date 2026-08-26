package domain

import "errors"

type ValidationIssue struct{ Field, Message string }
type ValidationReport struct{ Issues []ValidationIssue }

func (r ValidationReport) Valid() bool { return len(r.Issues) == 0 }
func (r ValidationReport) Error() error {
	if r.Valid() {
		return nil
	}
	return errors.New(r.Issues[0].Field + ": " + r.Issues[0].Message)
}
func ValidatePatient(p Patient) ValidationReport {
	r := ValidationReport{}
	if p.ID == "" {
		r.Issues = append(r.Issues, ValidationIssue{"id", "required"})
	}
	if p.Name == "" {
		r.Issues = append(r.Issues, ValidationIssue{"name", "required"})
	}
	return r
}
func ValidatePrescription(p Prescription) ValidationReport {
	r := ValidationReport{}
	if p.ID == "" {
		r.Issues = append(r.Issues, ValidationIssue{"id", "required"})
	}
	if p.PatientID == "" {
		r.Issues = append(r.Issues, ValidationIssue{"patient_id", "required"})
	}
	if p.Quantity < 1 {
		r.Issues = append(r.Issues, ValidationIssue{"quantity", "must be positive"})
	}
	return r
}
func ValidateTicket(t PickupTicket) ValidationReport {
	r := ValidationReport{}
	if t.Number < 1 {
		r.Issues = append(r.Issues, ValidationIssue{"number", "required"})
	}
	if t.Status == "" {
		r.Issues = append(r.Issues, ValidationIssue{"status", "required"})
	}
	return r
}
func ValidateDispense(d DispenseRecord) ValidationReport {
	r := ValidationReport{}
	if d.TicketID == "" {
		r.Issues = append(r.Issues, ValidationIssue{"ticket_id", "required"})
	}
	if d.Quantity < 1 {
		r.Issues = append(r.Issues, ValidationIssue{"quantity", "must be positive"})
	}
	return r
}
