# Suggested Price Calculator

This package implements the pricing logic for facility rentals, including special pricing rules based on a member's existing rentals.

## Overview

The `SuggestedPriceCalculator` calculates special prices for facility types based on what other facilities a member has already rented. This implements business rules like "If you have a box, canoes cost €80 and lockers cost €50" with absolute prices decided by the yacht club president.

## Architecture

- **Location**: Backend domain layer (`domain/facility_rental/pricing`)
- **Why Backend?**:
  - Business logic belongs in the domain
  - Prevents price manipulation from clients
  - Single source of truth for all clients (web, mobile, admin)
  - Easier to maintain and test

## Pricing Rules

Current pricing configuration (can be customized in `getDefaultPricingConfigs()`):

### Canoe (Facility Type ID: 2)

- **€80.00** if member has a Box (ID: 1)

### Locker (Facility Type ID: 3)

- **€50.00** if member has a Box (ID: 1)
- **€60.00** if member has a Mooring (ID: 4)

### How It Works

- If multiple rules apply, the **lowest price** is used
- Special prices are based on facilities the member has **already rented** in the same season
- Prices are absolute values set by the yacht club president, not percentages
- If no rule applies, the base suggested price is used

## Usage

### 1. In Domain Service

The `RentalManagementService` provides methods to calculate suggested prices:

```go
// Get suggested price with special pricing applied
suggestedPrice := rentalService.GetSuggestedPriceForMember(
    facilityTypeId,     // The facility type being rented
    baseSuggestedPrice, // The base price without special pricing
    memberId,           // The member requesting the rental
    seasonId,           // The season for the rental
)

// Get information about which pricing rules apply (for UI display)
applicablePricingRules := rentalService.GetApplicablePricingRulesForMember(
    facilityTypeId,
    memberId,
    seasonId,
)
```

### 2. Via HTTP API

**Endpoint**: `GET /api/v1.0/facilities/suggested-price`

**Query Parameters**:

- `facility_type_id` (required): ID of the facility type being rented
- `base_price` (required): The base suggested price
- `member_id` (required): ID of the member requesting the rental
- `season` (required): Season ID

**Example Request**:

```
GET /api/v1.0/facilities/suggested-price?facility_type_id=2&base_price=100.00&member_id=123&season=1
```

**Response**:

```json
{
  "suggestedPrice": 80.0,
  "basePrice": 100.0,
  "savingsAmount": 20.0,
  "hasSpecialPrice": true,
  "applicableRules": 1
}
```

### 3. In Frontend

```typescript
import { getSuggestedPrice } from "$lib/data/api/facilities-api";

const priceInfo = await getSuggestedPrice(
  facilityTypeId,
  basePrice,
  memberId,
  seasonId,
);

if (priceInfo.hasSpecialPrice) {
  console.log(`Special price available!`);
  console.log(`You save: €${priceInfo.savingsAmount}`);
  console.log(`Your price: €${priceInfo.suggestedPrice}`);
}
```

## Customizing Pricing Rules

### Option 1: Modify Default Configuration

Edit `getDefaultPricingConfigs()` in `suggested_price_calculator.go`:

```go
func getDefaultPricingConfigs() []FacilityTypePricingConfig {
    return []FacilityTypePricingConfig{
        // Add special pricing for new facility type
        {
            FacilityTypeId: 5, // Kayak
            PricingRules: []PricingRule{
                {
                    RequiredFacilityTypeId: 1,     // Box
                    SpecialPrice:           75.0,  // Fixed price in EUR
                },
            },
        },
    }
}
```

### Option 2: Load from Configuration (Future Enhancement)

The calculator supports dynamic configuration:

```go
calculator := pricing.NewSuggestedPriceCalculator()

// Load custom configuration
customConfigs := loadFromDatabase() // or config file
calculator.SetPricingConfigs(customConfigs)
```

## Testing

To test the calculator logic:

```go
calculator := pricing.NewSuggestedPriceCalculator()

// Member has rented a Box (ID: 1)
memberFacilityTypes := []int64{1}

// Calculate price for Canoe (ID: 2)
price := calculator.CalculateSuggestedPrice(
    2,     // Canoe
    100.0, // Base price
    memberFacilityTypes,
)
// price == 80.0 (special price of €80 applied)
```

## Important Notes

1. **Facility Type IDs**: The IDs in the configuration (1, 2, 3, 4) must match your actual database IDs. Verify these in your `facilities_catalog` table.

2. **Season-Based**: Discounts only apply based on facilities rented in the **same season**. Cross-season rentals don't provide discounts.

3. **Lowest Price Wins**: If multiple pricing rules apply, the lowest absolute price is used (prices don't stack).

4. **Cache Consideration**: When updating discount rules, consider clearing any cached pricing data in the frontend.

## Future Enhancements

- [ ] Store pricing rules in database
- [ ] Admin UI to configure special prices
- [ ] Time-based promotions (seasonal pricing)
- [ ] Member tier-based pricing
- [ ] Combination rules (multiple facilities required)
- [ ] Pricing rule expiration dates
- [ ] Minimum rental requirements for special prices
