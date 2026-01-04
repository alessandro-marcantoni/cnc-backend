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
	GetValidUntilDate() time.Time
}

type MembershipStatus string

const (
	MembershipStatusActive               MembershipStatus = "Active"
	MembershipStatusExpired              MembershipStatus = "Expired"
	MembershipStatusExclusionDeliberated MembershipStatus = "ExclusionDeliberated"
	MembershipStatusExcluded             MembershipStatus = "Excluded"
)

type Active struct {
	ValidUntilDate time.Time
}

type Expired struct {
	ValidUntilDate time.Time
}

type ExclusionDeliberated struct {
	ValidUntilDate time.Time
	DecisionDate   time.Time
}

type Excluded struct {
	ValidUntilDate time.Time
	DecisionDate   time.Time
}

func (a Active) GetStatus() MembershipStatus {
	return MembershipStatusActive
}

func (a Active) GetValidUntilDate() time.Time {
	return a.ValidUntilDate
}

func (e Expired) GetStatus() MembershipStatus {
	return MembershipStatusExpired
}

func (e Expired) GetValidUntilDate() time.Time {
	return e.ValidUntilDate
}

func (e ExclusionDeliberated) GetStatus() MembershipStatus {
	return MembershipStatusExclusionDeliberated
}

func (e ExclusionDeliberated) GetValidUntilDate() time.Time {
	return e.ValidUntilDate
}

func (e Excluded) GetStatus() MembershipStatus {
	return MembershipStatusExcluded
}

func (e Excluded) GetValidUntilDate() time.Time {
	return e.ValidUntilDate
}
