package payment

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type PaymentRepository interface {
	CreatePaymentForMembershipPeriod(membershipPeriodId int64, amount float64, currency string, paymentMethod string, transactionRef *string) result.Result[int64]
	CreatePaymentForRentedFacility(rentedFacilityId int64, amount float64, currency string, paymentMethod string, transactionRef *string) result.Result[int64]
	UpdatePayment(paymentId domain.Id[Payment], amount float64, currency string, paymentMethod string, transactionRef *string) result.Result[bool]
	DeletePayment(paymentId domain.Id[Payment]) result.Result[bool]
}
