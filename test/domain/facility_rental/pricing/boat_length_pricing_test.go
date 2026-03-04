package pricing

import (
	"math"
	"testing"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain/facility_rental/pricing"
)

func TestBoatLengthPriceCalculator_CalculatePriceForBoatLength(t *testing.T) {
	// Setup test tiers
	tiers := []pricing.BoatLengthTier{
		{
			MinLengthMeters: 0,
			MaxLengthMeters: 6.0,
			Price:           100.0,
		},
		{
			MinLengthMeters: 6.0,
			MaxLengthMeters: 8.0,
			Price:           130.0,
		},
		{
			MinLengthMeters: 8.0,
			MaxLengthMeters: 10.0,
			Price:           160.0,
		},
		{
			MinLengthMeters: 10.0,
			MaxLengthMeters: math.Inf(1),
			Price:           200.0,
		},
	}

	config := pricing.BoatLengthPricingConfig{
		FacilityTypeId: 1,
		Tiers:          tiers,
		DefaultPrice:   100.0,
	}

	calculator := pricing.NewBoatLengthPriceCalculator([]pricing.BoatLengthPricingConfig{config})

	tests := []struct {
		name             string
		facilityTypeId   int64
		boatLengthMeters float64
		expectedPrice    float64
	}{
		{
			name:             "Small boat in first tier",
			facilityTypeId:   1,
			boatLengthMeters: 5.5,
			expectedPrice:    100.0,
		},
		{
			name:             "Boat at exact tier boundary (lower)",
			facilityTypeId:   1,
			boatLengthMeters: 6.0,
			expectedPrice:    130.0,
		},
		{
			name:             "Medium boat in second tier",
			facilityTypeId:   1,
			boatLengthMeters: 7.0,
			expectedPrice:    130.0,
		},
		{
			name:             "Boat at exact tier boundary (upper)",
			facilityTypeId:   1,
			boatLengthMeters: 8.0,
			expectedPrice:    160.0,
		},
		{
			name:             "Large boat in third tier",
			facilityTypeId:   1,
			boatLengthMeters: 9.5,
			expectedPrice:    160.0,
		},
		{
			name:             "Very large boat in unlimited tier",
			facilityTypeId:   1,
			boatLengthMeters: 15.0,
			expectedPrice:    200.0,
		},
		{
			name:             "Huge boat in unlimited tier",
			facilityTypeId:   1,
			boatLengthMeters: 25.0,
			expectedPrice:    200.0,
		},
		{
			name:             "Zero length returns default",
			facilityTypeId:   1,
			boatLengthMeters: 0,
			expectedPrice:    100.0,
		},
		{
			name:             "Negative length returns default",
			facilityTypeId:   1,
			boatLengthMeters: -5.0,
			expectedPrice:    100.0,
		},
		{
			name:             "Non-configured facility type returns 0",
			facilityTypeId:   999,
			boatLengthMeters: 7.0,
			expectedPrice:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.CalculatePriceForBoatLength(tt.facilityTypeId, tt.boatLengthMeters)
			if result != tt.expectedPrice {
				t.Errorf("CalculatePriceForBoatLength() = %.2f, want %.2f", result, tt.expectedPrice)
			}
		})
	}
}

func TestBoatLengthPriceCalculator_HasBoatLengthPricing(t *testing.T) {
	config := pricing.BoatLengthPricingConfig{
		FacilityTypeId: 1,
		Tiers: []pricing.BoatLengthTier{
			{MinLengthMeters: 0, MaxLengthMeters: 10, Price: 100},
		},
		DefaultPrice: 100.0,
	}

	calculator := pricing.NewBoatLengthPriceCalculator([]pricing.BoatLengthPricingConfig{config})

	tests := []struct {
		name           string
		facilityTypeId int64
		expected       bool
	}{
		{
			name:           "Configured facility type",
			facilityTypeId: 1,
			expected:       true,
		},
		{
			name:           "Non-configured facility type",
			facilityTypeId: 999,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.HasBoatLengthPricing(tt.facilityTypeId)
			if result != tt.expected {
				t.Errorf("HasBoatLengthPricing() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBoatLengthPriceCalculator_GetPricingTiersForFacilityType(t *testing.T) {
	tiers := []pricing.BoatLengthTier{
		{MinLengthMeters: 0, MaxLengthMeters: 6, Price: 100},
		{MinLengthMeters: 6, MaxLengthMeters: 10, Price: 150},
	}

	config := pricing.BoatLengthPricingConfig{
		FacilityTypeId: 1,
		Tiers:          tiers,
		DefaultPrice:   100.0,
	}

	calculator := pricing.NewBoatLengthPriceCalculator([]pricing.BoatLengthPricingConfig{config})

	t.Run("Get tiers for configured facility", func(t *testing.T) {
		result, exists := calculator.GetPricingTiersForFacilityType(1)
		if !exists {
			t.Error("Expected facility type to exist")
		}
		if len(result) != 2 {
			t.Errorf("Expected 2 tiers, got %d", len(result))
		}
	})

	t.Run("Get tiers for non-configured facility", func(t *testing.T) {
		result, exists := calculator.GetPricingTiersForFacilityType(999)
		if exists {
			t.Error("Expected facility type to not exist")
		}
		if result != nil {
			t.Error("Expected nil result for non-existent facility")
		}
	})
}

func TestBoatLengthPriceCalculator_GetDefaultPrice(t *testing.T) {
	config := pricing.BoatLengthPricingConfig{
		FacilityTypeId: 1,
		Tiers: []pricing.BoatLengthTier{
			{MinLengthMeters: 0, MaxLengthMeters: 10, Price: 100},
		},
		DefaultPrice: 75.0,
	}

	calculator := pricing.NewBoatLengthPriceCalculator([]pricing.BoatLengthPricingConfig{config})

	t.Run("Get default price for configured facility", func(t *testing.T) {
		result, exists := calculator.GetDefaultPrice(1)
		if !exists {
			t.Error("Expected facility type to exist")
		}
		if result != 75.0 {
			t.Errorf("Expected default price 75.0, got %.2f", result)
		}
	})

	t.Run("Get default price for non-configured facility", func(t *testing.T) {
		result, exists := calculator.GetDefaultPrice(999)
		if exists {
			t.Error("Expected facility type to not exist")
		}
		if result != 0 {
			t.Errorf("Expected 0 for non-existent facility, got %.2f", result)
		}
	})
}

func TestValidateTiers(t *testing.T) {
	tests := []struct {
		name        string
		tiers       []pricing.BoatLengthTier
		expectError bool
	}{
		{
			name:        "Empty tiers is valid",
			tiers:       []pricing.BoatLengthTier{},
			expectError: false,
		},
		{
			name: "Valid continuous tiers",
			tiers: []pricing.BoatLengthTier{
				{MinLengthMeters: 0, MaxLengthMeters: 5, Price: 100},
				{MinLengthMeters: 5, MaxLengthMeters: 10, Price: 150},
				{MinLengthMeters: 10, MaxLengthMeters: math.Inf(1), Price: 200},
			},
			expectError: false,
		},
		{
			name: "Valid tiers with gaps",
			tiers: []pricing.BoatLengthTier{
				{MinLengthMeters: 0, MaxLengthMeters: 5, Price: 100},
				{MinLengthMeters: 7, MaxLengthMeters: 10, Price: 150},
			},
			expectError: false,
		},
		{
			name: "Invalid overlapping tiers",
			tiers: []pricing.BoatLengthTier{
				{MinLengthMeters: 0, MaxLengthMeters: 6, Price: 100},
				{MinLengthMeters: 5, MaxLengthMeters: 10, Price: 150},
			},
			expectError: true,
		},
		{
			name: "Invalid tier - max less than min",
			tiers: []pricing.BoatLengthTier{
				{MinLengthMeters: 10, MaxLengthMeters: 5, Price: 100},
			},
			expectError: true,
		},
		{
			name: "Invalid tier - max equals min",
			tiers: []pricing.BoatLengthTier{
				{MinLengthMeters: 5, MaxLengthMeters: 5, Price: 100},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidateTiers(tt.tiers)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestCompositePriceCalculator_CalculatePrice(t *testing.T) {
	// Setup discount calculator
	discountConfigs := []pricing.FacilityTypePricingConfig{
		{
			FacilityTypeId: 1,
			PricingRules: []pricing.PricingRule{
				{
					RequiredFacilityTypeId: 2, // If member has facility type 2
					SpecialPrice:           80.0,
				},
			},
		},
	}
	discountCalc := pricing.NewSuggestedPriceCalculator(discountConfigs)

	// Setup boat length calculator
	boatLengthConfigs := []pricing.BoatLengthPricingConfig{
		{
			FacilityTypeId: 1,
			Tiers: []pricing.BoatLengthTier{
				{MinLengthMeters: 0, MaxLengthMeters: 6, Price: 100},
				{MinLengthMeters: 6, MaxLengthMeters: 10, Price: 130},
				{MinLengthMeters: 10, MaxLengthMeters: math.Inf(1), Price: 160},
			},
			DefaultPrice: 100.0,
		},
	}
	boatLengthCalc := pricing.NewBoatLengthPriceCalculator(boatLengthConfigs)

	// Create composite calculator
	composite := pricing.NewCompositePriceCalculator(discountCalc, boatLengthCalc)

	tests := []struct {
		name                      string
		ctx                       pricing.PriceCalculationContext
		expectedPrice             float64
		expectedMethod            pricing.PricingMethod
		expectedBoatLengthApplied bool
		expectedDiscountApplied   bool
	}{
		{
			name: "Base pricing - no boat length, no discount",
			ctx: pricing.PriceCalculationContext{
				FacilityTypeId:            1,
				BaseSuggestedPrice:        100.0,
				MemberRentedFacilityTypes: []int64{},
				BoatLengthMeters:          nil,
			},
			expectedPrice:             100.0,
			expectedMethod:            pricing.BasePricing,
			expectedBoatLengthApplied: false,
			expectedDiscountApplied:   false,
		},
		{
			name: "Boat length pricing only",
			ctx: pricing.PriceCalculationContext{
				FacilityTypeId:            1,
				BaseSuggestedPrice:        100.0,
				MemberRentedFacilityTypes: []int64{},
				BoatLengthMeters:          floatPtr(7.0),
			},
			expectedPrice:             130.0,
			expectedMethod:            pricing.BoatLengthPricing,
			expectedBoatLengthApplied: true,
			expectedDiscountApplied:   false,
		},
		{
			name: "Discount pricing only",
			ctx: pricing.PriceCalculationContext{
				FacilityTypeId:            1,
				BaseSuggestedPrice:        100.0,
				MemberRentedFacilityTypes: []int64{2}, // Has facility type 2
				BoatLengthMeters:          nil,
			},
			expectedPrice:             80.0,
			expectedMethod:            pricing.DiscountPricing,
			expectedBoatLengthApplied: false,
			expectedDiscountApplied:   true,
		},
		{
			name: "Combined pricing - boat length + discount",
			ctx: pricing.PriceCalculationContext{
				FacilityTypeId:            1,
				BaseSuggestedPrice:        100.0,
				MemberRentedFacilityTypes: []int64{2}, // Has facility type 2
				BoatLengthMeters:          floatPtr(7.0),
			},
			expectedPrice:             80.0, // Discount applied to boat length price
			expectedMethod:            pricing.CombinedPricing,
			expectedBoatLengthApplied: true,
			expectedDiscountApplied:   true,
		},
		{
			name: "Large boat with discount",
			ctx: pricing.PriceCalculationContext{
				FacilityTypeId:            1,
				BaseSuggestedPrice:        100.0,
				MemberRentedFacilityTypes: []int64{2},
				BoatLengthMeters:          floatPtr(12.0),
			},
			expectedPrice:             80.0, // Discount trumps boat length
			expectedMethod:            pricing.CombinedPricing,
			expectedBoatLengthApplied: true,
			expectedDiscountApplied:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := composite.CalculatePrice(tt.ctx)

			if result.FinalPrice != tt.expectedPrice {
				t.Errorf("FinalPrice = %.2f, want %.2f", result.FinalPrice, tt.expectedPrice)
			}
			if result.PricingMethod != tt.expectedMethod {
				t.Errorf("PricingMethod = %s, want %s", result.PricingMethod, tt.expectedMethod)
			}
			if result.BoatLengthTierApplied != tt.expectedBoatLengthApplied {
				t.Errorf("BoatLengthTierApplied = %v, want %v", result.BoatLengthTierApplied, tt.expectedBoatLengthApplied)
			}
			if result.DiscountApplied != tt.expectedDiscountApplied {
				t.Errorf("DiscountApplied = %v, want %v", result.DiscountApplied, tt.expectedDiscountApplied)
			}
		})
	}
}

func TestCompositePriceCalculator_CalculateSimplePrice(t *testing.T) {
	discountCalc := pricing.NewSuggestedPriceCalculator([]pricing.FacilityTypePricingConfig{})
	boatLengthCalc := pricing.NewBoatLengthPriceCalculator([]pricing.BoatLengthPricingConfig{})
	composite := pricing.NewCompositePriceCalculator(discountCalc, boatLengthCalc)

	ctx := pricing.PriceCalculationContext{
		FacilityTypeId:            1,
		BaseSuggestedPrice:        100.0,
		MemberRentedFacilityTypes: []int64{},
		BoatLengthMeters:          nil,
	}

	result := composite.CalculateSimplePrice(ctx)
	if result != 100.0 {
		t.Errorf("CalculateSimplePrice() = %.2f, want 100.00", result)
	}
}

// Helper function
func floatPtr(f float64) *float64 {
	return &f
}
