package pricing

// PricingRule represents a special price that applies when a member already has a certain facility type
type PricingRule struct {
	// RequiredFacilityTypeId is the facility type the member must already have
	RequiredFacilityTypeId int64
	// SpecialPrice is the absolute price to apply (e.g., 80.00 EUR)
	SpecialPrice float64
}

// FacilityTypePricingConfig holds all pricing rules for a specific facility type
type FacilityTypePricingConfig struct {
	FacilityTypeId int64
	PricingRules   []PricingRule
}

// SuggestedPriceCalculator calculates suggested prices based on member's existing rentals
type SuggestedPriceCalculator struct {
	pricingConfigs []FacilityTypePricingConfig
}

// NewSuggestedPriceCalculator creates a new calculator with the provided pricing rules
func NewSuggestedPriceCalculator(pricingConfigs []FacilityTypePricingConfig) *SuggestedPriceCalculator {
	return &SuggestedPriceCalculator{
		pricingConfigs: pricingConfigs,
	}
}

// CalculateSuggestedPrice calculates the suggested price for a facility type based on member's existing rentals
// Returns the special price if a rule applies, otherwise returns the base suggested price
func (c *SuggestedPriceCalculator) CalculateSuggestedPrice(
	facilityTypeId int64,
	baseSuggestedPrice float64,
	memberRentedFacilityTypes []int64,
) float64 {
	// Find pricing config for this facility type
	var config *FacilityTypePricingConfig
	for i := range c.pricingConfigs {
		if c.pricingConfigs[i].FacilityTypeId == facilityTypeId {
			config = &c.pricingConfigs[i]
			break
		}
	}

	// No pricing rules for this facility type
	if config == nil {
		return baseSuggestedPrice
	}

	// Find the best price that applies (lowest price wins)
	bestPrice := baseSuggestedPrice
	priceFound := false

	for _, rule := range config.PricingRules {
		// Check if member has the required facility type
		if c.hasFacilityType(memberRentedFacilityTypes, rule.RequiredFacilityTypeId) {
			if !priceFound || rule.SpecialPrice < bestPrice {
				bestPrice = rule.SpecialPrice
				priceFound = true
			}
		}
	}

	return bestPrice
}

// hasFacilityType checks if a facility type is in the list
func (c *SuggestedPriceCalculator) hasFacilityType(
	rentedTypes []int64,
	targetType int64,
) bool {
	for _, t := range rentedTypes {
		if t == targetType {
			return true
		}
	}
	return false
}

// GetApplicablePricingRules returns all pricing rules that apply to a member for a given facility type
// This is useful for displaying pricing information in the UI
func (c *SuggestedPriceCalculator) GetApplicablePricingRules(
	facilityTypeId int64,
	memberRentedFacilityTypes []int64,
) []PricingRule {
	// Find pricing config for this facility type
	var config *FacilityTypePricingConfig
	for i := range c.pricingConfigs {
		if c.pricingConfigs[i].FacilityTypeId == facilityTypeId {
			config = &c.pricingConfigs[i]
			break
		}
	}

	if config == nil {
		return []PricingRule{}
	}

	// Filter rules that apply
	applicableRules := []PricingRule{}
	for _, rule := range config.PricingRules {
		if c.hasFacilityType(memberRentedFacilityTypes, rule.RequiredFacilityTypeId) {
			applicableRules = append(applicableRules, rule)
		}
	}

	return applicableRules
}

// GetAllPricingConfigs returns all pricing configurations (useful for admin UI)
func (c *SuggestedPriceCalculator) GetAllPricingConfigs() []FacilityTypePricingConfig {
	return c.pricingConfigs
}

// SetPricingConfigs allows updating the pricing configuration
// This could be used to load configuration from a database or config file
func (c *SuggestedPriceCalculator) SetPricingConfigs(configs []FacilityTypePricingConfig) {
	c.pricingConfigs = configs
}
