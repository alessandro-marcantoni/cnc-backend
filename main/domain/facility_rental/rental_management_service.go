package facilityrental

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/facility_rental/pricing"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type RentalManagementService struct {
	repository            FacilityRepository
	waitingListRepository WaitingListRepository
	priceCalculator       *pricing.SuggestedPriceCalculator
}

func NewRentalManagementService(repository FacilityRepository, waitingListRepository WaitingListRepository) *RentalManagementService {
	// Fetch pricing rules from database
	pricingRules := repository.GetPricingRules()

	// Build pricing configs from repository data
	pricingConfigs := buildPricingConfigs(pricingRules)

	return &RentalManagementService{
		repository:            repository,
		waitingListRepository: waitingListRepository,
		priceCalculator:       pricing.NewSuggestedPriceCalculator(pricingConfigs),
	}
}

// buildPricingConfigs converts repository pricing rules into pricing calculator configs
func buildPricingConfigs(rules []PricingRule) []pricing.FacilityTypePricingConfig {
	// Group rules by facility type
	configMap := make(map[int64][]pricing.PricingRule)

	for _, rule := range rules {
		if !rule.Active {
			continue
		}

		facilityTypeId := rule.FacilityTypeId.Value

		pricingRule := pricing.PricingRule{
			RequiredFacilityTypeId: rule.RequiredFacilityTypeId.Value,
			SpecialPrice:           rule.SpecialPrice,
		}

		configMap[facilityTypeId] = append(configMap[facilityTypeId], pricingRule)
	}

	// Convert map to slice
	configs := make([]pricing.FacilityTypePricingConfig, 0, len(configMap))
	for facilityTypeId, rules := range configMap {
		configs = append(configs, pricing.FacilityTypePricingConfig{
			FacilityTypeId: facilityTypeId,
			PricingRules:   rules,
		})
	}

	return configs
}

func (this RentalManagementService) RentService(
	facilityId domain.Id[Facility],
	memberId domain.Id[membership.User],
	season int64,
	price float64,
	boat *BoatInfo,
) result.Result[RentedFacility] {
	// Rent the facility
	rentResult := this.repository.RentFacility(memberId, facilityId, season, price, boat)
	if !rentResult.IsSuccess() {
		return rentResult
	}

	// Get the rented facility to access its type
	rentedFacility := rentResult.Value()
	facility := rentedFacility.GetFacility()
	facilityTypeId := facility.FacilityType.Id

	// Remove member from waiting list if they were waiting for this facility type
	// Convert User ID to Member ID (they share the same underlying value)
	memberIdForWaitlist := domain.Id[membership.Member]{Value: memberId.Value}

	// Attempt to remove - if they weren't on the list, it's not an error
	_ = this.waitingListRepository.RemoveEntryByMemberAndType(facilityTypeId, memberIdForWaitlist)
	// We ignore the result because:
	// - If successful, great! Member was waiting and is now removed
	// - If they weren't on the waiting list, that's fine too
	// - Either way, the facility rental succeeded

	return rentResult
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

// GetSuggestedPriceForMember calculates the suggested price for a facility type
// considering any discounts based on the member's existing rentals
func (this RentalManagementService) GetSuggestedPriceForMember(
	facilityTypeId domain.Id[FacilityType],
	baseSuggestedPrice float64,
	memberId domain.Id[membership.User],
	season int64,
) float64 {
	// Get member's currently rented facilities for the season
	rentedFacilities := this.repository.GetFacilitiesRentedByMember(memberId, season)

	// Extract the facility type IDs as int64
	rentedFacilityTypeIds := make([]int64, len(rentedFacilities))
	for i, rentedFacility := range rentedFacilities {
		facility := rentedFacility.GetFacility()
		rentedFacilityTypeIds[i] = facility.FacilityType.Id.Value
	}

	// Calculate suggested price with applicable discounts
	return this.priceCalculator.CalculateSuggestedPrice(
		facilityTypeId.Value,
		baseSuggestedPrice,
		rentedFacilityTypeIds,
	)
}

// GetApplicableDiscountsForMember returns all pricing rules that apply to a member
// for a given facility type, useful for displaying in the UI
func (this RentalManagementService) GetApplicableDiscountsForMember(
	facilityTypeId domain.Id[FacilityType],
	memberId domain.Id[membership.User],
	season int64,
) []pricing.PricingRule {
	// Get member's currently rented facilities for the season
	rentedFacilities := this.repository.GetFacilitiesRentedByMember(memberId, season)

	// Extract the facility type IDs as int64
	rentedFacilityTypeIds := make([]int64, len(rentedFacilities))
	for i, rentedFacility := range rentedFacilities {
		facility := rentedFacility.GetFacility()
		rentedFacilityTypeIds[i] = facility.FacilityType.Id.Value
	}

	// Get applicable pricing rules
	return this.priceCalculator.GetApplicablePricingRules(
		facilityTypeId.Value,
		rentedFacilityTypeIds,
	)
}

func (this RentalManagementService) FreeFacility(rentedFacilityId domain.Id[RentedFacility]) result.Result[bool] {
	return this.repository.FreeFacility(rentedFacilityId)
}
