package membership

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type MemberManagementService struct {
	repository MemberRepository
}

func NewMemberManagementService(repository MemberRepository) *MemberManagementService {
	return &MemberManagementService{repository: repository}
}

func (this MemberManagementService) GetListOfAllMembers() result.Result[[]Member] {
	return this.repository.GetAllMembers()
}

func (this MemberManagementService) GetListOfMembersBySeason(year int64) result.Result[[]Member] {
	return this.repository.GetMembersBySeason(year)
}

func (this MemberManagementService) GetMemberById(id domain.Id[Member], season string) result.Result[MemberDetails] {
	return this.repository.GetMemberById(id, season)
}

func (this MemberManagementService) GetMembersWhoDidNotPayForServices() []Member {
	return this.repository.GetMembersWhoDidNotPayForServices()
}

func (this MemberManagementService) GetMembersWhoDidNotPayForMembership() []Member {
	return this.repository.GetMembersWhoDidNotPayForMembership()
}

func (this MemberManagementService) CreateMember(user User, createMembership bool, seasonId *int64, price *float64) result.Result[MemberDetails] {
	return this.repository.CreateMember(user, createMembership, seasonId, price)
}

func (this MemberManagementService) AddMembership(memberId domain.Id[Member], seasonId int64, seasonStartsAt string, seasonEndsAt string, price float64) result.Result[MemberDetails] {
	return this.repository.AddMembership(memberId, seasonId, seasonStartsAt, seasonEndsAt, price)
}
