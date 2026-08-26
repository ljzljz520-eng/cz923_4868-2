package service

import "errors"

func EnsurePharmacist(name string) error {
	if name == "" {
		return errors.New("pharmacist required")
	}
	return nil
}
func EnsureQuantity(q int) error {
	if q < 1 {
		return errors.New("quantity must be positive")
	}
	if q > 100 {
		return errors.New("quantity exceeds counter limit")
	}
	return nil
}
func CanCancel(status string) bool { return status != "completed" && status != "cancelled" }
func CanCall(status string) bool   { return status == "waiting" }
