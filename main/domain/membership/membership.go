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
	Price   float64
	Payment payment.Payment
}

type MembershipInfo interface {
	GetStatus() MembershipStatus
	GetValidFromDate() time.Time
	GetValidUntilDate() time.Time
	GetPeriodId() *int64
}

type MembershipStatus string

const (
	MembershipStatusActive    MembershipStatus = "ACTIVE"
	MembershipStatusSuspended MembershipStatus = "SUSPENDED"
	MembershipStatusExcluded  MembershipStatus = "EXCLUDED"
	MembershipStatusExpired   MembershipStatus = "EXPIRED"
	MembershipStatusNone      MembershipStatus = "NONE"
)

type Active struct {
	PeriodId       int64
	ValidFromDate  time.Time
	ValidUntilDate time.Time
}

type Suspended struct {
	PeriodId       int64
	ValidFromDate  time.Time
	ValidUntilDate time.Time
}

type Excluded struct {
	PeriodId       int64
	ValidFromDate  time.Time
	ValidUntilDate time.Time
	ExcludedAt     time.Time
}

type Expired struct {
	PeriodId       int64
	ValidFromDate  time.Time
	ValidUntilDate time.Time
}

type None struct {
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

func (a Active) GetPeriodId() *int64 {
	return &a.PeriodId
}

func (s Suspended) GetStatus() MembershipStatus {
	return MembershipStatusSuspended
}

func (s Suspended) GetValidFromDate() time.Time {
	return s.ValidFromDate
}

func (s Suspended) GetValidUntilDate() time.Time {
	return s.ValidUntilDate
}

func (s Suspended) GetPeriodId() *int64 {
	return &s.PeriodId
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

func (e Excluded) GetPeriodId() *int64 {
	return &e.PeriodId
}

func (e Expired) GetStatus() MembershipStatus {
	return MembershipStatusExpired
}

func (e Expired) GetValidFromDate() time.Time {
	return e.ValidFromDate
}

func (e Expired) GetValidUntilDate() time.Time {
	return e.ValidUntilDate
}

func (e Expired) GetPeriodId() *int64 {
	return &e.PeriodId
}

func (n None) GetStatus() MembershipStatus {
	return MembershipStatusNone
}

func (n None) GetValidFromDate() time.Time {
	return time.Time{}
}

func (n None) GetValidUntilDate() time.Time {
	return time.Time{}
}

func (n None) GetPeriodId() *int64 {
	return nil
}
