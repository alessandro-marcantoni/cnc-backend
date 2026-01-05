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

func (this MemberManagementService) GetUpdatedListOfMembers() result.Result[[]Member] {
	return this.repository.GetAllMembers()
}

func (this MemberManagementService) GetMemberById(id domain.Id[Member]) result.Result[MemberDetails] {
	return this.repository.GetMemberById(id)
}

func (this MemberManagementService) GetMembersWhoDidNotPayForServices() []Member {
	return this.repository.GetMembersWhoDidNotPayForServices()
}

func (this MemberManagementService) GetMembersWhoDidNotPayForMembership() []Member {
	return this.repository.GetMembersWhoDidNotPayForMembership()
}
