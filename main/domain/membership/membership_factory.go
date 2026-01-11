package membership

import (
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain/payment"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

func RenewedMembership(currentMembership Membership) result.Result[Membership] {
	newValidityDate := currentMembership.Status.GetValidUntilDate().AddDate(1, 0, 0)

	return result.Ok(Membership{
		Number: currentMembership.Number,
		Status: Active{
			ValidUntilDate: newValidityDate,
		},
		Payment: payment.PaymentUnpaid{},
	})
}

func ExcludedMembership(currentMembership Membership, decisionDate time.Time) Membership {
	return Membership{
		Number: currentMembership.Number,
		Status: Excluded{
			ValidFromDate:  currentMembership.Status.GetValidFromDate(),
			ValidUntilDate: currentMembership.Status.GetValidUntilDate(),
			ExcludedAt:     decisionDate,
		},
		Payment: currentMembership.Payment,
	}
}
