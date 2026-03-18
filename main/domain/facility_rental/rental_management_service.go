package facilityrental

import (
	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/facility_rental/pricing"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type RentalManagementService struct {
	repository               FacilityRepository
	waitingListRepository    WaitingListRepository
	priceCalculator          *pricing.SuggestedPriceCalculator
	compositePriceCalculator *pricing.CompositePriceCalculator
}

func NewRentalManagementService(repository FacilityRepository, waitingListRepository WaitingListRepository) *RentalManagementService {
	// Fetch pricing rules from database
	pricingRules := repository.GetPricingRules()

	// Build pricing configs from repository data
	pricingConfigs := buildPricingConfigs(pricingRules)

	// Create discount-based price calculator
	discountCalculator := pricing.NewSuggestedPriceCalculator(pricingConfigs)

	// Build boat length pricing configs
	boatLengthConfigs := buildBoatLengthPricingConfigs(repository)

	// Create boat-length price calculator
	boatLengthCalculator := pricing.NewBoatLengthPriceCalculator(boatLengthConfigs)

	// Create composite price calculator
	compositePriceCalculator := pricing.NewCompositePriceCalculator(discountCalculator, boatLengthCalculator)

	return &RentalManagementService{
		repository:               repository,
		waitingListRepository:    waitingListRepository,
		priceCalculator:          discountCalculator,
		compositePriceCalculator: compositePriceCalculator,
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

// buildBoatLengthPricingConfigs creates boat-length pricing configurations from database
func buildBoatLengthPricingConfigs(repository FacilityRepository) []pricing.BoatLengthPricingConfig {
	// Load boat length pricing tiers from database
	dbTiers := repository.GetBoatLengthPricingTiers()

	// Group tiers by facility type
	configMap := make(map[int64][]pricing.BoatLengthTier)

	for _, dbTier := range dbTiers {
		facilityTypeId := dbTier.FacilityTypeId.Value

		// Convert database tier to pricing tier
		pricingTier := pricing.BoatLengthTier{
			MinLengthMeters: dbTier.MinLengthMeters,
			MaxLengthMeters: 0, // Will be set below
			Price:           dbTier.Price,
		}

		// Handle max length (nil means infinity)
		if dbTier.MaxLengthMeters == nil {
			// Use a very large number to represent infinity for practical purposes
			pricingTier.MaxLengthMeters = 1e9
		} else {
			pricingTier.MaxLengthMeters = *dbTier.MaxLengthMeters
		}

		configMap[facilityTypeId] = append(configMap[facilityTypeId], pricingTier)
	}

	// Convert map to slice of configs
	configs := make([]pricing.BoatLengthPricingConfig, 0, len(configMap))
	catalog := repository.GetFacilitiesCatalog()

	for facilityTypeId, tiers := range configMap {
		// Find the facility type to get default price
		var defaultPrice float64 = 0
		for _, facilityType := range catalog {
			if facilityType.Id.Value == facilityTypeId {
				defaultPrice = facilityType.SuggestedPrice
				break
			}
		}

		configs = append(configs, pricing.BoatLengthPricingConfig{
			FacilityTypeId: facilityTypeId,
			Tiers:          tiers,
			DefaultPrice:   defaultPrice,
		})
	}

	return configs
}

func (this RentalManagementService) RentService(
	facilityId domain.Id[Facility],
	memberId domain.Id[membership.User],
	season int64,
	price float64,
	discountApplied bool,
	boat *BoatInfo,
	leerboard *LeerboardInfo,
) result.Result[RentedFacility] {
	// Rent the facility
	rentResult := this.repository.RentFacility(memberId, facilityId, season, price, discountApplied, boat, leerboard)
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

// GetSuggestedPriceWithBoatLength calculates the suggested price for a facility
// considering boat length (if applicable) and discounts based on member's existing rentals
func (this RentalManagementService) GetSuggestedPriceWithBoatLength(
	facilityTypeId domain.Id[FacilityType],
	baseSuggestedPrice float64,
	memberId domain.Id[membership.User],
	season int64,
	boatLengthMeters *float64,
) pricing.PriceCalculationResult {
	// Get member's currently rented facilities for the season
	rentedFacilities := this.repository.GetFacilitiesRentedByMember(memberId, season)

	// Extract the facility type IDs as int64
	// Also check if member already has a rental with discount applied
	rentedFacilityTypeIds := make([]int64, len(rentedFacilities))
	memberHasDiscountedRental := false
	for i, rentedFacility := range rentedFacilities {
		facility := rentedFacility.GetFacility()
		rentedFacilityTypeIds[i] = facility.FacilityType.Id.Value

		// Check if any existing rental has discount applied
		if rentedFacility.GetDiscountApplied() {
			memberHasDiscountedRental = true
		}
	}

	// Create pricing context
	ctx := pricing.PriceCalculationContext{
		FacilityTypeId:            facilityTypeId.Value,
		BaseSuggestedPrice:        baseSuggestedPrice,
		MemberRentedFacilityTypes: rentedFacilityTypeIds,
		MemberHasDiscountedRental: memberHasDiscountedRental,
		BoatLengthMeters:          boatLengthMeters,
	}

	// Calculate price using composite calculator
	return this.compositePriceCalculator.CalculatePrice(ctx)
}

// GetBoatLengthTiers returns the pricing tiers for a boat facility
func (this RentalManagementService) GetBoatLengthTiers(
	facilityTypeId domain.Id[FacilityType],
) ([]pricing.BoatLengthTier, bool) {
	return this.compositePriceCalculator.GetPricingInformation(
			facilityTypeId.Value,
			0, // Base price not needed for just getting tiers
			nil,
		).BoatLengthTiers, this.compositePriceCalculator.GetPricingInformation(
			facilityTypeId.Value,
			0,
			nil,
		).HasBoatLengthPricing
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
