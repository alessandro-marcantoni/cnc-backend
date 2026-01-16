package facilityrental

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type RentalManagementService struct {
	repository FacilityRepository
}

func NewRentalManagementService(repository FacilityRepository) *RentalManagementService {
	return &RentalManagementService{repository: repository}
}

func (this RentalManagementService) RentService(
	facilityId domain.Id[Facility],
	memberId domain.Id[membership.User],
	season int64,
	price float64,
	boat *BoatInfo,
) result.Result[RentedFacility] {
	return this.repository.RentFacility(memberId, facilityId, season, price, boat)
}

func (this RentalManagementService) GetFacilitiesCatalog() []FacilityType {
	return this.repository.GetFacilitiesCatalog()
}

func (this RentalManagementService) GetFacilitiesByType(facilityTypeId domain.Id[FacilityType], seasonId int64) []FacilityWithStatus {
	return this.repository.GetFacilitiesByType(facilityTypeId, seasonId)
}

func (this RentalManagementService) GetFacilitiesRentedByMember(memberId domain.Id[membership.User], season int64) []RentedFacility {
	return this.repository.GetFacilitiesRentedByMember(memberId, season)
}
