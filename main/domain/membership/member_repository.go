package membership

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type MemberRepository interface {
	GetAllMembers() []Member
	GetMemberById(id domain.Id[Member]) result.Result[Member]
	GetMembersWhoDidNotPayForServices() []Member
	GetMembersWhoDidNotPayForMembership() []Member
}
