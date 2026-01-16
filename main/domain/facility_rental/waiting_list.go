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
	AddEntry(entry WaitingListEntry) result.Result[WaitingListEntry]
	RemoveEntry(entryId domain.Id[WaitingListEntry]) result.Result[WaitingListEntry]
	RemoveEntryByMemberAndType(facilityType domain.Id[FacilityType], memberId domain.Id[membership.Member]) result.Result[WaitingListEntry]
	GetMemberEntry(facilityType domain.Id[FacilityType], memberId domain.Id[membership.Member]) result.Result[WaitingListEntry]
}

type WaitingList struct {
	FacilityType FacilityType
	Entries      []WaitingListEntry
}

type WaitingListEntry struct {
	Id           domain.Id[WaitingListEntry]
	MemberId     domain.Id[membership.Member]
	FacilityType domain.Id[FacilityType]
	QueuedAt     time.Time
	Notes        string
}

func NewWaitingListEntry(
	memberId domain.Id[membership.Member],
	facilityType domain.Id[FacilityType],
	notes string,
) WaitingListEntry {
	return WaitingListEntry{
		Id:           domain.Id[WaitingListEntry]{Value: 0}, // ID will be set by database
		MemberId:     memberId,
		FacilityType: facilityType,
		QueuedAt:     time.Now(),
		Notes:        notes,
	}
}

func ReconstructWaitingListEntry(
	id domain.Id[WaitingListEntry],
	memberId domain.Id[membership.Member],
	facilityType domain.Id[FacilityType],
	queuedAt time.Time,
	notes string,
) WaitingListEntry {
	return WaitingListEntry{
		Id:           id,
		MemberId:     memberId,
		FacilityType: facilityType,
		QueuedAt:     queuedAt,
		Notes:        notes,
	}
}
