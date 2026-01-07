package membership

import (
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/payment"
)

const SuggestedMembershipPrice float64 = 130.0

type Membership struct {
	Id      domain.Id[Membership]
	Number  int64
	Status  MembershipInfo
	Payment payment.Payment
}

type MembershipInfo interface {
	GetStatus() MembershipStatus
	GetValidFromDate() time.Time
	GetValidUntilDate() time.Time
}

type MembershipStatus string

const (
	MembershipStatusActive   MembershipStatus = "ACTIVE"
	MembershipStatusInactive MembershipStatus = "INACTIVE"
)

type Active struct {
	ValidFromDate  time.Time
	ValidUntilDate time.Time
}

type Inactive struct {
	ValidFromDate  time.Time
	ValidUntilDate time.Time
	ExcludedAt     time.Time
}

func (a Active) GetStatus() MembershipStatus {
	return MembershipStatusActive
}

func (a Active) GetValidFromDate() time.Time {
	return a.ValidFromDate
}

func (a Active) GetValidUntilDate() time.Time {
	return a.ValidUntilDate
}

func (e Inactive) GetStatus() MembershipStatus {
	return MembershipStatusInactive
}

func (e Inactive) GetValidFromDate() time.Time {
	return e.ValidFromDate
}

func (e Inactive) GetValidUntilDate() time.Time {
	return e.ValidUntilDate
}
