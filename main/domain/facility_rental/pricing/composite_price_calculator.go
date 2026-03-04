package pricing

// CompositePriceCalculator combines multiple pricing strategies
// to calculate the final suggested price for a facility rental
type CompositePriceCalculator struct {
	discountCalculator   *SuggestedPriceCalculator
	boatLengthCalculator *BoatLengthPriceCalculator
}

// NewCompositePriceCalculator creates a new composite calculator
func NewCompositePriceCalculator(
	discountCalculator *SuggestedPriceCalculator,
	boatLengthCalculator *BoatLengthPriceCalculator,
) *CompositePriceCalculator {
	return &CompositePriceCalculator{
		discountCalculator:   discountCalculator,
		boatLengthCalculator: boatLengthCalculator,
	}
}

// PriceCalculationContext holds all information needed to calculate a price
type PriceCalculationContext struct {
	FacilityTypeId            int64
	BaseSuggestedPrice        float64
	MemberRentedFacilityTypes []int64
	BoatLengthMeters          *float64 // Optional: only for boat facilities
}

// PriceCalculationResult holds the result of price calculation with details
type PriceCalculationResult struct {
	FinalPrice            float64
	BasePrice             float64
	PricingMethod         PricingMethod
	DiscountApplied       bool
	DiscountAmount        float64
	BoatLengthTierApplied bool
	BoatLengthTierPrice   float64
}

// PricingMethod indicates which pricing strategy was used
type PricingMethod string

const (
	BasePricing       PricingMethod = "BASE"        // No special pricing applied
	DiscountPricing   PricingMethod = "DISCOUNT"    // Discount based on owned facilities
	BoatLengthPricing PricingMethod = "BOAT_LENGTH" // Price based on boat length
	CombinedPricing   PricingMethod = "COMBINED"    // Both discount and boat length applied
)

// CalculatePrice calculates the final price using all available pricing strategies
// Priority: Boat length pricing > Discount pricing > Base price
func (c *CompositePriceCalculator) CalculatePrice(ctx PriceCalculationContext) PriceCalculationResult {
	result := PriceCalculationResult{
		BasePrice:     ctx.BaseSuggestedPrice,
		PricingMethod: BasePricing,
	}

	// Strategy 1: Check if boat-length-based pricing applies
	hasBoatLengthPricing := c.boatLengthCalculator.HasBoatLengthPricing(ctx.FacilityTypeId)

	if hasBoatLengthPricing && ctx.BoatLengthMeters != nil && *ctx.BoatLengthMeters > 0 {
		// Boat length pricing takes precedence
		boatLengthPrice := c.boatLengthCalculator.CalculatePriceForBoatLength(
			ctx.FacilityTypeId,
			*ctx.BoatLengthMeters,
		)

		if boatLengthPrice > 0 {
			result.BoatLengthTierApplied = true
			result.BoatLengthTierPrice = boatLengthPrice
			result.FinalPrice = boatLengthPrice
			result.PricingMethod = BoatLengthPricing

			// Check if discount can be applied on top of boat length price
			discountPrice := c.discountCalculator.CalculateSuggestedPrice(
				ctx.FacilityTypeId,
				boatLengthPrice, // Use boat length price as base
				ctx.MemberRentedFacilityTypes,
			)

			if discountPrice < boatLengthPrice {
				result.DiscountApplied = true
				result.DiscountAmount = boatLengthPrice - discountPrice
				result.FinalPrice = discountPrice
				result.PricingMethod = CombinedPricing
			}

			return result
		}
	}

	// Strategy 2: Apply discount-based pricing
	discountPrice := c.discountCalculator.CalculateSuggestedPrice(
		ctx.FacilityTypeId,
		ctx.BaseSuggestedPrice,
		ctx.MemberRentedFacilityTypes,
	)

	if discountPrice < ctx.BaseSuggestedPrice {
		result.DiscountApplied = true
		result.DiscountAmount = ctx.BaseSuggestedPrice - discountPrice
		result.FinalPrice = discountPrice
		result.PricingMethod = DiscountPricing
	} else {
		// No special pricing applies, use base price
		result.FinalPrice = ctx.BaseSuggestedPrice
	}

	return result
}

// CalculateSimplePrice is a convenience method that returns just the final price
func (c *CompositePriceCalculator) CalculateSimplePrice(ctx PriceCalculationContext) float64 {
	result := c.CalculatePrice(ctx)
	return result.FinalPrice
}

// GetPricingInformation returns detailed pricing information for UI display
func (c *CompositePriceCalculator) GetPricingInformation(
	facilityTypeId int64,
	baseSuggestedPrice float64,
	memberRentedFacilityTypes []int64,
) PricingInformation {
	info := PricingInformation{
		FacilityTypeId: facilityTypeId,
		BasePrice:      baseSuggestedPrice,
	}

	// Check for boat length pricing
	if tiers, hasBoatLengthPricing := c.boatLengthCalculator.GetPricingTiersForFacilityType(facilityTypeId); hasBoatLengthPricing {
		info.HasBoatLengthPricing = true
		info.BoatLengthTiers = tiers
	}

	// Get applicable discount rules
	discountRules := c.discountCalculator.GetApplicablePricingRules(
		facilityTypeId,
		memberRentedFacilityTypes,
	)

	if len(discountRules) > 0 {
		info.HasDiscounts = true
		info.ApplicableDiscounts = discountRules
	}

	return info
}

// PricingInformation holds all pricing information for a facility type
type PricingInformation struct {
	FacilityTypeId       int64
	BasePrice            float64
	HasBoatLengthPricing bool
	BoatLengthTiers      []BoatLengthTier
	HasDiscounts         bool
	ApplicableDiscounts  []PricingRule
}
