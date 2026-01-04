package payment

import (
	"time"
)

type Payment interface {
	GetStatus() PaymentStatus
	GetAmount() float64
	GetPaidAt() time.Time
}

type PaymentPaid struct {
	AmountPaid  float64
	PaymentDate time.Time
}

func (p PaymentPaid) GetStatus() PaymentStatus {
	return Paid
}

func (p PaymentPaid) GetAmount() float64 {
	return p.AmountPaid
}

func (p PaymentPaid) GetPaidAt() time.Time {
	return p.PaymentDate
}

type PaymentUnpaid struct {
	AmountDue float64
	DueDate   time.Time
}

func (p PaymentUnpaid) GetStatus() PaymentStatus {
	return Unpaid
}

func (p PaymentUnpaid) GetAmount() float64 {
	return p.AmountDue
}

func (p PaymentUnpaid) GetPaidAt() time.Time {
	return time.Time{} // Return zero time for unpaid payments
}
