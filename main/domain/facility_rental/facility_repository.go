package facilityrental

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type FacilityRepository interface {
	GetFacilitiesCatalog() []FacilityType
	GetFacilitiesByType(facilityTypeId domain.Id[FacilityType], seasonId int64) []FacilityWithStatus
	GetFacilityById(facilityId domain.Id[Facility]) (FacilityWithStatus, bool)
	GetAvailableFacilities(serviceType FacilityName) []Facility
	GetFacilitiesRentedByMember(memberId domain.Id[membership.User], season int64) []RentedFacility
	GetPricingRules() []PricingRule
	GetBoatLengthPricingTiers() []BoatLengthPricingTier
	RentFacility(
		memberId domain.Id[membership.User],
		facilityId domain.Id[Facility],
		season int64,
		price float64,
		discountApplied bool,
		boatInfo *BoatInfo,
		leerboardInfo *LeerboardInfo,
	) result.Result[RentedFacility]
	ChangeFacility(rentedFacilityId domain.Id[RentedFacility], newFacilityId domain.Id[Facility]) result.Result[RentedFacility]
	FreeFacility(rentedFacilityId domain.Id[RentedFacility]) result.Result[bool]
}

type PricingRule struct {
	Id                     domain.Id[PricingRule]
	FacilityTypeId         domain.Id[FacilityType]
	RequiredFacilityTypeId domain.Id[FacilityType]
	SpecialPrice           float64
	Description            string
	Active                 bool
}

type BoatLengthPricingTier struct {
	Id              domain.Id[BoatLengthPricingTier]
	FacilityTypeId  domain.Id[FacilityType]
	MinLengthMeters float64
	MaxLengthMeters *float64 // nil means no upper limit (infinity)
	Price           float64
	Currency        string
	Active          bool
}
