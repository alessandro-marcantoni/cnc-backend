package payment

import (
	"time"
)

type Payment interface {
	GetStatus() PaymentStatus
}

type PaymentPaid struct {
	ID             int64
	AmountPaid     float64
	Currency       string
	PaymentDate    time.Time
	PaymentMethod  string
	TransactionRef string
}

func (p PaymentPaid) GetStatus() PaymentStatus {
	return Paid
}

type PaymentUnpaid struct {
}

func (p PaymentUnpaid) GetStatus() PaymentStatus {
	return Unpaid
}
