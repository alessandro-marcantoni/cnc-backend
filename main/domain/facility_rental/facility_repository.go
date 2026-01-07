package facilityrental

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type FacilityRepository interface {
	GetFacilitiesCatalog() []FacilityType
	GetFacilitiesByType(facilityTypeId domain.Id[FacilityType]) []FacilityWithStatus
	GetAvailableFacilities(serviceType FacilityName) []Facility
	GetFacilitiesRentedByMember(memberId domain.Id[membership.User]) []RentedFacility
	RentFacility(memberId domain.Id[membership.User], facilityId domain.Id[Facility], boatInfo *domain.Id[BoatInfo]) result.Result[RentedFacility]
}
