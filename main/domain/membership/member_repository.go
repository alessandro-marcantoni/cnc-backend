package membership

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type MemberRepository interface {
	GetAllMembers() result.Result[[]Member]
	GetMemberById(id domain.Id[Member], season string) result.Result[MemberDetails]
	GetMembersBySeason(year int64) result.Result[[]Member]
	GetMembersWhoDidNotPayForServices() []Member
	GetMembersWhoDidNotPayForMembership() []Member
	CreateMember(user User, createMembership bool, seasonId *int64, price *float64) result.Result[MemberDetails]
	AddMembership(memberId domain.Id[Member], seasonId int64, seasonStartsAt string, seasonEndsAt string, price float64) result.Result[MemberDetails]
}
