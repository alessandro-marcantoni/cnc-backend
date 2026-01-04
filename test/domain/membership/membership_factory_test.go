package membership_test

import (
	"testing"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/payment"
)

func TestRenewedMembership(t *testing.T) {
	// Setup initial membership
	currentDate := domain.NewDate(2025, 1, 1).Value()
	memberId := domain.NewId[membership.Membership](1)
	currentMembership := membership.Membership{
		Number: memberId,
		Status: membership.Active{
			ValidUntilDate: currentDate,
		},
		Payment: payment.PaymentPaid{
			AmountPaid:  membership.SuggestedMembershipPrice,
			PaymentDate: currentDate,
		},
	}

	// Test renewal
	renewedResult := membership.RenewedMembership(currentMembership)

	if !renewedResult.IsSuccess() {
		t.Errorf("Expected successful renewal, got error: %v", renewedResult.Error())
		return
	}

	renewed := renewedResult.Value()

	// Check membership number remains the same
	if renewed.Number != memberId {
		t.Errorf("Expected membership number to be %v, got %v", memberId, renewed.Number)
	}

	// Check validity date is extended by 1 year
	expectedNewDate := domain.NewDate(2026, 1, 1).Value()
	if renewed.Status.GetValidUntilDate() != expectedNewDate {
		t.Errorf("Expected new validity date to be %v, got %v", expectedNewDate, renewed.Status.GetValidUntilDate())
	}

	// Check payment is set to unpaid with correct amount
	payment, ok := renewed.Payment.(payment.PaymentUnpaid)
	if !ok {
		t.Error("Expected payment to be unpaid")
		return
	}
	if payment.AmountDue != membership.SuggestedMembershipPrice {
		t.Errorf("Expected amount due to be %v, got %v", membership.SuggestedMembershipPrice, payment.AmountDue)
	}
}

func TestDeliberatedExclusionMembership(t *testing.T) {
	// Setup initial membership
	validUntilDate := domain.NewDate(2025, 1, 31).Value()
	decisionDate := domain.NewDate(2025, 6, 15).Value()
	memberId := domain.NewId[membership.Membership](1)
	currentMembership := membership.Membership{
		Number: memberId,
		Status: membership.Active{
			ValidUntilDate: validUntilDate,
		},
		Payment: payment.PaymentUnpaid{
			AmountDue: membership.SuggestedMembershipPrice,
			DueDate:   validUntilDate,
		},
	}

	// Test deliberated exclusion
	excluded := membership.DeliberatedExclusionMembership(currentMembership, decisionDate)

	// Check membership number remains the same
	if excluded.Number != memberId {
		t.Errorf("Expected membership number to be %v, got %v", memberId, excluded.Number)
	}

	// Check the status is changed to ExclusionDeliberated
	status, ok := excluded.Status.(membership.ExclusionDeliberated)
	if !ok {
		t.Error("Expected status to be ExclusionDeliberated")
		return
	}

	// Check dates are set correctly
	if status.ValidUntilDate != validUntilDate {
		t.Errorf("Expected valid until date to be %v, got %v", validUntilDate, status.ValidUntilDate)
	}
	if status.DecisionDate != decisionDate {
		t.Errorf("Expected decision date to be %v, got %v", decisionDate, status.DecisionDate)
	}
}

func TestExcludedMembership(t *testing.T) {
	// Setup initial membership
	validUntilDate := domain.NewDate(2025, 1, 31).Value()
	decisionDate := domain.NewDate(2025, 6, 15).Value()
	memberId := domain.NewId[membership.Membership](1)
	currentMembership := membership.Membership{
		Number: memberId,
		Status: membership.Active{
			ValidUntilDate: validUntilDate,
		},
		Payment: payment.PaymentUnpaid{
			AmountDue: membership.SuggestedMembershipPrice,
			DueDate:   validUntilDate,
		},
	}

	// Test exclusion
	excluded := membership.ExcludedMembership(currentMembership, decisionDate)

	// Check membership number remains the same
	if excluded.Number != memberId {
		t.Errorf("Expected membership number to be %v, got %v", memberId, excluded.Number)
	}

	// Check the status is changed to Excluded
	status, ok := excluded.Status.(membership.Excluded)
	if !ok {
		t.Error("Expected status to be Excluded")
		return
	}

	// Check dates are set correctly
	if status.ValidUntilDate != validUntilDate {
		t.Errorf("Expected valid until date to be %v, got %v", validUntilDate, status.ValidUntilDate)
	}
	if status.DecisionDate != decisionDate {
		t.Errorf("Expected decision date to be %v, got %v", decisionDate, status.DecisionDate)
	}
}
