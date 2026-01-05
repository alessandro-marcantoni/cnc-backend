package membership

import (
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/payment"
)

const SuggestedMembershipPrice float64 = 130.0

type Membership struct {
	Number  domain.Id[Membership]
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
	MembershipStatusActive               MembershipStatus = "ACTIVE"
	MembershipStatusUnpaid               MembershipStatus = "UNPAID"
	MembershipStatusExclusionDeliberated MembershipStatus = "EXCLUSION_DELIBERATED"
	MembershipStatusExcluded             MembershipStatus = "EXCLUDED"
)

type Active struct {
	ValidFromDate  time.Time
	ValidUntilDate time.Time
}

type Unpaid struct {
	ValidFromDate  time.Time
	ValidUntilDate time.Time
}

type ExclusionDeliberated struct {
	ValidFromDate  time.Time
	ValidUntilDate time.Time
	DecisionDate   time.Time
}

type Excluded struct {
	ValidFromDate  time.Time
	ValidUntilDate time.Time
	DecisionDate   time.Time
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

func (e Unpaid) GetStatus() MembershipStatus {
	return MembershipStatusUnpaid
}

func (e Unpaid) GetValidFromDate() time.Time {
	return e.ValidFromDate
}

func (e Unpaid) GetValidUntilDate() time.Time {
	return e.ValidUntilDate
}

func (e ExclusionDeliberated) GetStatus() MembershipStatus {
	return MembershipStatusExclusionDeliberated
}

func (e ExclusionDeliberated) GetValidFromDate() time.Time {
	return e.ValidFromDate
}

func (e ExclusionDeliberated) GetValidUntilDate() time.Time {
	return e.ValidUntilDate
}

func (e Excluded) GetStatus() MembershipStatus {
	return MembershipStatusExcluded
}

func (e Excluded) GetValidFromDate() time.Time {
	return e.ValidFromDate
}

func (e Excluded) GetValidUntilDate() time.Time {
	return e.ValidUntilDate
}
