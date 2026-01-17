package membership

import (
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/shared/errors"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type Member struct {
	User
	Membership          Membership
	HasUnpaidFacilities bool
}

type MemberDetails struct {
	User
	Memberships []Membership
}

func (m Member) IsActive() bool {
	return m.Membership.Status.GetStatus() == MembershipStatusActive
}

func (m Member) CanRentServices() bool {
	return m.IsActive()
}

func (m Member) RenewMembership() result.Result[Member] {
	if m.Membership.Status.GetStatus() != MembershipStatusActive {
		return result.Err[Member](errors.MembershipStatusError{Description: "only active members can renew their membership"})
	}
	return result.Map(RenewedMembership(m.Membership), func(newMembership Membership) Member {
		return Member{
			User:       m.User,
			Membership: newMembership,
		}
	})
}

func (m Member) Exclude(decisionDate time.Time) result.Result[Member] {
	if m.Membership.Status.GetStatus() != MembershipStatusActive {
		return result.Err[Member](errors.MembershipStatusError{Description: "only members who are active can be excluded"})
	}
	return result.Ok(Member{
		User:       m.User,
		Membership: ExcludedMembership(m.Membership, decisionDate),
	})
}
