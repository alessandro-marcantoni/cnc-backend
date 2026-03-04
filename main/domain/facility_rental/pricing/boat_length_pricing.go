package pricing

import "math"

// BoatLengthTier represents a price tier based on boat length
type BoatLengthTier struct {
	// MinLengthMeters is the minimum boat length (inclusive) for this tier
	MinLengthMeters float64
	// MaxLengthMeters is the maximum boat length (exclusive) for this tier
	// Use math.Inf(1) for no upper limit
	MaxLengthMeters float64
	// Price is the price for boats in this length range
	Price float64
}

// BoatLengthPricingConfig holds the configuration for boat-length-based pricing
type BoatLengthPricingConfig struct {
	FacilityTypeId int64
	Tiers          []BoatLengthTier
	// DefaultPrice is used when no tier matches or if boat length is not provided
	DefaultPrice float64
}

// BoatLengthPriceCalculator calculates prices based on boat length
type BoatLengthPriceCalculator struct {
	configs map[int64]BoatLengthPricingConfig
}

// NewBoatLengthPriceCalculator creates a new calculator with the provided configurations
func NewBoatLengthPriceCalculator(configs []BoatLengthPricingConfig) *BoatLengthPriceCalculator {
	configMap := make(map[int64]BoatLengthPricingConfig)
	for _, config := range configs {
		configMap[config.FacilityTypeId] = config
	}

	return &BoatLengthPriceCalculator{
		configs: configMap,
	}
}

// CalculatePriceForBoatLength calculates the price for a facility based on boat length
// Returns the tier price if a matching tier is found, otherwise returns the default price
func (c *BoatLengthPriceCalculator) CalculatePriceForBoatLength(
	facilityTypeId int64,
	boatLengthMeters float64,
) float64 {
	// Find pricing config for this facility type
	config, exists := c.configs[facilityTypeId]
	if !exists {
		// No boat-length pricing configured for this facility type
		return 0 // Caller should use base suggested price
	}

	// If boat length is invalid (zero or negative), return default
	if boatLengthMeters <= 0 {
		return config.DefaultPrice
	}

	// Find matching tier
	for _, tier := range config.Tiers {
		if boatLengthMeters >= tier.MinLengthMeters && boatLengthMeters < tier.MaxLengthMeters {
			return tier.Price
		}
	}

	// No matching tier found, return default
	return config.DefaultPrice
}

// GetPricingTiersForFacilityType returns all pricing tiers for a facility type
// Useful for displaying pricing information in the UI
func (c *BoatLengthPriceCalculator) GetPricingTiersForFacilityType(
	facilityTypeId int64,
) ([]BoatLengthTier, bool) {
	config, exists := c.configs[facilityTypeId]
	if !exists {
		return nil, false
	}

	return config.Tiers, true
}

// HasBoatLengthPricing checks if a facility type has boat-length-based pricing configured
func (c *BoatLengthPriceCalculator) HasBoatLengthPricing(facilityTypeId int64) bool {
	_, exists := c.configs[facilityTypeId]
	return exists
}

// GetDefaultPrice returns the default price for a facility type when no tier matches
func (c *BoatLengthPriceCalculator) GetDefaultPrice(facilityTypeId int64) (float64, bool) {
	config, exists := c.configs[facilityTypeId]
	if !exists {
		return 0, false
	}

	return config.DefaultPrice, true
}

// ValidateTiers checks if tiers are properly configured (no gaps or overlaps)
func ValidateTiers(tiers []BoatLengthTier) error {
	if len(tiers) == 0 {
		return nil // Empty is valid
	}

	// Sort tiers by min length (assuming they're provided in order)
	// Check for gaps and overlaps
	for i := 0; i < len(tiers)-1; i++ {
		currentTier := tiers[i]
		nextTier := tiers[i+1]

		// Check if current tier's max equals next tier's min (no gaps/overlaps)
		if currentTier.MaxLengthMeters != nextTier.MinLengthMeters {
			// There's either a gap or an overlap
			if currentTier.MaxLengthMeters > nextTier.MinLengthMeters {
				// Overlap
				return &ValidationError{
					Message:   "tiers overlap",
					TierIndex: i,
				}
			}
			// Gap is allowed, as default price will be used
		}

		// Check if max > min
		if currentTier.MaxLengthMeters <= currentTier.MinLengthMeters {
			return &ValidationError{
				Message:   "tier max length must be greater than min length",
				TierIndex: i,
			}
		}
	}

	// Validate last tier
	lastTier := tiers[len(tiers)-1]
	if lastTier.MaxLengthMeters <= lastTier.MinLengthMeters && !math.IsInf(lastTier.MaxLengthMeters, 1) {
		return &ValidationError{
			Message:   "tier max length must be greater than min length or infinity",
			TierIndex: len(tiers) - 1,
		}
	}

	return nil
}

// ValidationError represents a tier validation error
type ValidationError struct {
	Message   string
	TierIndex int
}

func (e *ValidationError) Error() string {
	return e.Message
}
