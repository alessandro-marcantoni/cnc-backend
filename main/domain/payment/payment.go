package payment

import (
	"time"
)

type Payment interface {
	GetStatus() PaymentStatus
}

type PaymentPaid struct {
	AmountPaid  float64
	Currency    string
	PaymentDate time.Time
}

func (p PaymentPaid) GetStatus() PaymentStatus {
	return Paid
}

type PaymentUnpaid struct {
}

func (p PaymentUnpaid) GetStatus() PaymentStatus {
	return Unpaid
}
