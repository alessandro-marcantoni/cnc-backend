package payment

type PaymentStatus string

const (
	Paid   PaymentStatus = "PAID"
	Unpaid PaymentStatus = "UNPAID"
)
