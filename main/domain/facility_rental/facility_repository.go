package facilityrental

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type FacilityRepository interface {
	GetFacilitiesCatalog() []FacilityType
	GetFacilitiesByType(facilityTypeId domain.Id[FacilityType], seasonId int64) []FacilityWithStatus
	GetAvailableFacilities(serviceType FacilityName) []Facility
	GetFacilitiesRentedByMember(memberId domain.Id[membership.User], season int64) []RentedFacility
	GetPricingRules() []PricingRule
	RentFacility(
		memberId domain.Id[membership.User],
		facilityId domain.Id[Facility],
		season int64,
		price float64,
		boatInfo *BoatInfo,
	) result.Result[RentedFacility]
}

type PricingRule struct {
	Id                     domain.Id[PricingRule]
	FacilityTypeId         domain.Id[FacilityType]
	RequiredFacilityTypeId domain.Id[FacilityType]
	SpecialPrice           float64
	Description            string
	Active                 bool
}
