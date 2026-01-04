package facilityrental

import (
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type WaitingListRepository interface {
	GetWaitingList(facilityType domain.Id[FacilityType]) result.Result[WaitingList]
	GetNextEntry(facilityType domain.Id[FacilityType]) result.Result[WaitingListEntry]
	AddEntry(facilityType domain.Id[FacilityType], memberId domain.Id[membership.Member]) result.Result[WaitingListEntry]
	RemoveEntry(facilityType domain.Id[FacilityType], memberId domain.Id[membership.Member]) result.Result[WaitingListEntry]
}

type WaitingList struct {
	FacilityType FacilityType
	Entries      []WaitingListEntry
}

type WaitingListEntry struct {
	MemberId     domain.Id[membership.Member]
	WaitingSince time.Time
}
