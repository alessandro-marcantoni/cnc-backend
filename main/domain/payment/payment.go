package payment

import (
	"time"
)

type Payment interface {
	GetStatus() PaymentStatus
}

type PaymentPaid struct {
	AmountPaid  float64
	PaymentDate time.Time
}

func (p PaymentPaid) GetStatus() PaymentStatus {
	return Paid
}

type PaymentUnpaid struct {
	AmountDue float64
	DueDate   time.Time
}

func (p PaymentUnpaid) GetStatus() PaymentStatus {
	return Unpaid
}
