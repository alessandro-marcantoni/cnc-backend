package membership

type MemberManagementService struct {
	repository MemberRepository
}

func (this MemberManagementService) GetUpdatedListOfMembers() []Member {
	return this.repository.GetAllMembers()
}

func (this MemberManagementService) GetMembersWhoDidNotPayForServices() []Member {
	return this.repository.GetMembersWhoDidNotPayForServices()
}

func (this MemberManagementService) GetMembersWhoDidNotPayForMembership() []Member {
	return this.repository.GetMembersWhoDidNotPayForMembership()
}

