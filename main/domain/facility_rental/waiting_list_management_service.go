package facilityrental

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type WaitingListManagementService struct {
	repository WaitingListRepository
}

func NewWaitingListManagementService(repository WaitingListRepository) *WaitingListManagementService {
	return &WaitingListManagementService{repository: repository}
}

// AddToWaitingList adds a member to the waiting list for a facility type
func (s *WaitingListManagementService) AddToWaitingList(
	memberId domain.Id[membership.Member],
	facilityTypeId domain.Id[FacilityType],
	notes string,
) result.Result[WaitingListEntry] {
	// Check if member is already in the waiting list
	existingEntry := s.repository.GetMemberEntry(facilityTypeId, memberId)
	if existingEntry.IsSuccess() {
		// Member already in waiting list, return the existing entry
		return existingEntry
	}

	// Create new entry
	entry := NewWaitingListEntry(memberId, facilityTypeId, notes)
	return s.repository.AddEntry(entry)
}

// RemoveFromWaitingList removes a member from the waiting list
func (s *WaitingListManagementService) RemoveFromWaitingList(
	entryId domain.Id[WaitingListEntry],
) result.Result[WaitingListEntry] {
	return s.repository.RemoveEntry(entryId)
}

// RemoveFromWaitingListByMemberAndType removes a member from the waiting list by member ID and facility type
func (s *WaitingListManagementService) RemoveFromWaitingListByMemberAndType(
	memberId domain.Id[membership.Member],
	facilityTypeId domain.Id[FacilityType],
) result.Result[WaitingListEntry] {
	return s.repository.RemoveEntryByMemberAndType(facilityTypeId, memberId)
}

// GetWaitingList returns the waiting list for a facility type, ordered by queued_at
func (s *WaitingListManagementService) GetWaitingList(
	facilityTypeId domain.Id[FacilityType],
) result.Result[WaitingList] {
	return s.repository.GetWaitingList(facilityTypeId)
}

// GetNextInLine returns the next member in the waiting list for a facility type
func (s *WaitingListManagementService) GetNextInLine(
	facilityTypeId domain.Id[FacilityType],
) result.Result[WaitingListEntry] {
	return s.repository.GetNextEntry(facilityTypeId)
}

// GetMemberEntry checks if a member is in the waiting list for a facility type
func (s *WaitingListManagementService) GetMemberEntry(
	memberId domain.Id[membership.Member],
	facilityTypeId domain.Id[FacilityType],
) result.Result[WaitingListEntry] {
	return s.repository.GetMemberEntry(facilityTypeId, memberId)
}
