# Boat Length-Based Pricing

## Overview

The boat-length-based pricing feature allows facilities that require boats to have dynamic pricing based on the length of the boat being stored. This provides a fair pricing model where larger boats (which require more space and maintenance) cost more than smaller boats.

## Architecture

### Components

1. **BoatLengthPriceCalculator** (`pricing/boat_length_pricing.go`)
   - Calculates prices based on boat length tiers
   - Supports multiple tiers with min/max length ranges
   - Falls back to default price if no tier matches

2. **CompositePriceCalculator** (`pricing/composite_price_calculator.go`)
   - Combines boat-length pricing with discount-based pricing
   - Determines which pricing strategy to apply
   - Returns detailed pricing information

3. **RentalManagementService** (`rental_management_service.go`)
   - Orchestrates pricing calculations
   - Configures boat-length tiers for facilities with boats
   - Exposes pricing methods to handlers

## Pricing Tiers

### Default Configuration

The system is configured with the following default tiers for boat facilities:

| Boat Length Range | Price Multiplier | Example (Base: €100) |
|-------------------|------------------|----------------------|
| 0m - 6m          | 1.0x             | €100                |
| 6m - 8m          | 1.3x             | €130                |
| 8m - 10m         | 1.6x             | €160                |
| 10m+             | 2.0x             | €200                |

### Customization

Currently, tiers are hardcoded in `buildBoatLengthPricingConfigs()`. Future enhancements could:
- Store tiers in the database
- Allow per-facility-type tier configuration
- Add admin UI for tier management

## Pricing Strategy Priority

The system applies pricing strategies in the following order:

1. **Boat Length Pricing** (if facility has boat and length provided)
   - Determines base price from boat length tier
   - Can be combined with discounts

2. **Discount Pricing** (if member has qualifying facilities)
   - Applied to base price or boat-length price
   - Based on facility pricing rules

3. **Base Price** (fallback)
   - Uses facility type's suggested price

### Pricing Method Types

```go
type PricingMethod string

const (
    BasePricing       PricingMethod = "BASE"        // No special pricing
    DiscountPricing   PricingMethod = "DISCOUNT"    // Discount applied
    BoatLengthPricing PricingMethod = "BOAT_LENGTH" // Boat length tier
    CombinedPricing   PricingMethod = "COMBINED"    // Both strategies
)
```

## API Usage

### Get Suggested Price with Boat Length

**Endpoint:** `GET /api/v1.0/suggested-price`

**Query Parameters:**

- `facility_type_id` (required): ID of the facility type
- `member_id` (required): ID of the member
- `season` (optional): Season ID (defaults to 0)
- `boat_length` (optional): Boat length in meters

**Example Request:**

```bash
curl "http://localhost:8080/api/v1.0/suggested-price?facility_type_id=2&member_id=5&season=1&boat_length=7.5"
```

**Response:**

```json
{
  "suggestedPrice": 130.0,
  "basePrice": 100.0,
  "pricingMethod": "BOAT_LENGTH",
  "discountApplied": false,
  "discountAmount": 0,
  "boatLengthTierApplied": true,
  "boatLengthTierPrice": 130.0,
  "applicableRules": 0,
  "boatLengthTiers": [
    {
      "minLengthMeters": 0,
      "maxLengthMeters": 6.0,
      "price": 100.0
    },
    {
      "minLengthMeters": 6.0,
      "maxLengthMeters": 8.0,
      "price": 130.0
    },
    {
      "minLengthMeters": 8.0,
      "maxLengthMeters": 10.0,
      "price": 160.0
    },
    {
      "minLengthMeters": 10.0,
      "maxLengthMeters": null,
      "price": 200.0
    }
  ],
  "hasBoatLengthPricing": true
}
```

### Combined Pricing Example

If a member with boat length 7.5m also has a discount (e.g., owns a Box):

**Response:**

```json
{
  "suggestedPrice": 104.0,
  "basePrice": 100.0,
  "pricingMethod": "COMBINED",
  "discountApplied": true,
  "discountAmount": 26.0,
  "boatLengthTierApplied": true,
  "boatLengthTierPrice": 130.0,
  "applicableRules": 1,
  "boatLengthTiers": [...],
  "hasBoatLengthPricing": true
}
```

In this case:

1. Boat length 7.5m → €130 (tier price)
2. Member has Box → 20% discount applied to €130
3. Final price: €104

## Code Examples

### Using the Composite Price Calculator

```go
// Create pricing context
ctx := pricing.PriceCalculationContext{
    FacilityTypeId:            facilityTypeId.Value,
    BaseSuggestedPrice:        100.0,
    MemberRentedFacilityTypes: []int64{1, 3}, // Member has Box and Locker
    BoatLengthMeters:          &boatLength,   // 7.5 meters
}

// Calculate price
result := compositePriceCalculator.CalculatePrice(ctx)

// Access results
fmt.Printf("Final Price: €%.2f\n", result.FinalPrice)
fmt.Printf("Pricing Method: %s\n", result.PricingMethod)
fmt.Printf("Boat Tier Applied: %v\n", result.BoatLengthTierApplied)
fmt.Printf("Discount Applied: %v\n", result.DiscountApplied)
```

### Defining Custom Tiers

```go
tiers := []pricing.BoatLengthTier{
    {
        MinLengthMeters: 0,
        MaxLengthMeters: 5.0,
        Price:           80.0,
    },
    {
        MinLengthMeters: 5.0,
        MaxLengthMeters: 7.5,
        Price:           120.0,
    },
    {
        MinLengthMeters: 7.5,
        MaxLengthMeters: math.Inf(1), // No upper limit
        Price:           180.0,
    },
}

config := pricing.BoatLengthPricingConfig{
    FacilityTypeId: 2,
    Tiers:          tiers,
    DefaultPrice:   100.0,
}
```

## Frontend Integration

### Display Boat Length Tiers

```typescript
// Fetch pricing information
const response = await fetch(
  `/api/v1.0/suggested-price?facility_type_id=${facilityTypeId}&member_id=${memberId}&boat_length=${boatLength}`,
);
const data = await response.json();

// Show tiers to user
if (data.hasBoatLengthPricing) {
  console.log("Boat Length Pricing Tiers:");
  data.boatLengthTiers.forEach((tier) => {
    console.log(
      `${tier.minLengthMeters}m - ${tier.maxLengthMeters}m: €${tier.price}`,
    );
  });
}

// Display final price with breakdown
console.log(`Base Price: €${data.basePrice}`);
if (data.boatLengthTierApplied) {
  console.log(`Boat Length Tier: €${data.boatLengthTierPrice}`);
}
if (data.discountApplied) {
  console.log(`Discount: -€${data.discountAmount}`);
}
console.log(`Final Price: €${data.suggestedPrice}`);
```

### Dynamic Boat Length Input

```svelte
<script>
  let boatLength = 7.5;
  let pricing = { suggestedPrice: 0, boatLengthTiers: [] };

  async function updatePrice() {
    const response = await fetch(
      `/api/v1.0/suggested-price?facility_type_id=${facilityTypeId}&member_id=${memberId}&boat_length=${boatLength}`
    );
    pricing = await response.json();
  }

  $: boatLength, updatePrice(); // Update when boat length changes
</script>

<label>
  Boat Length (meters):
  <input type="number" bind:value={boatLength} step="0.1" min="0" />
</label>

<div class="price-display">
  <h3>Suggested Price: €{pricing.suggestedPrice.toFixed(2)}</h3>

  {#if pricing.hasBoatLengthPricing}
    <h4>Pricing Tiers:</h4>
    <ul>
      {#each pricing.boatLengthTiers as tier}
        <li class:active={boatLength >= tier.minLengthMeters && boatLength < tier.maxLengthMeters}>
          {tier.minLengthMeters}m - {tier.maxLengthMeters === null ? '∞' : tier.maxLengthMeters + 'm'}:
          €{tier.price.toFixed(2)}
        </li>
      {/each}
    </ul>
  {/if}
</div>
```

## Database Schema (Future Enhancement)

To persist boat length pricing tiers in the database:

```sql
CREATE TABLE boat_length_pricing_tiers (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    facility_type_id BIGINT NOT NULL REFERENCES facilities_catalog(id),
    min_length_meters NUMERIC(10,2) NOT NULL,
    max_length_meters NUMERIC(10,2),  -- NULL for infinity
    price NUMERIC(10,2) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_length_range CHECK (
        max_length_meters IS NULL OR max_length_meters > min_length_meters
    )
);

CREATE INDEX idx_boat_pricing_facility_type
    ON boat_length_pricing_tiers(facility_type_id);
```

## Testing

### Unit Tests

```go
func TestBoatLengthPriceCalculator(t *testing.T) {
    tiers := []pricing.BoatLengthTier{
        {MinLengthMeters: 0, MaxLengthMeters: 6, Price: 100},
        {MinLengthMeters: 6, MaxLengthMeters: 8, Price: 130},
        {MinLengthMeters: 8, MaxLengthMeters: math.Inf(1), Price: 160},
    }

    config := pricing.BoatLengthPricingConfig{
        FacilityTypeId: 1,
        Tiers:         tiers,
        DefaultPrice:  100,
    }

    calc := pricing.NewBoatLengthPriceCalculator([]pricing.BoatLengthPricingConfig{config})

    // Test tier matching
    assert.Equal(t, 100.0, calc.CalculatePriceForBoatLength(1, 5.5))
    assert.Equal(t, 130.0, calc.CalculatePriceForBoatLength(1, 7.0))
    assert.Equal(t, 160.0, calc.CalculatePriceForBoatLength(1, 12.0))
}
```

### Integration Tests

```bash
# Test boat length pricing
curl "http://localhost:8080/api/v1.0/suggested-price?facility_type_id=2&member_id=1&boat_length=5"
# Expected: Base tier price

curl "http://localhost:8080/api/v1.0/suggested-price?facility_type_id=2&member_id=1&boat_length=15"
# Expected: Highest tier price

# Test combined pricing (boat length + discount)
# First, rent a facility that provides discount
curl -X POST http://localhost:8080/api/v1.0/rented-facilities \
  -H "Content-Type: application/json" \
  -d '{"facilityId": 1, "memberId": 1, "seasonId": 1, "price": 50}'

# Then check price with boat length
curl "http://localhost:8080/api/v1.0/suggested-price?facility_type_id=2&member_id=1&boat_length=7&season=1"
# Expected: Boat tier price with discount applied
```

## Troubleshooting

### Issue: Boat length pricing not applied

**Symptoms:**

- API returns base price even with boat_length parameter
- `hasBoatLengthPricing` is false

**Solutions:**

1. Verify facility type has `HasBoat = true`
2. Check boat length is positive and > 0
3. Ensure tiers are configured in `buildBoatLengthPricingConfigs()`

### Issue: Wrong tier selected

**Symptoms:**

- Price doesn't match expected tier
- Boat length 6.0m gets wrong tier

**Solutions:**

1. Check tier boundaries (MaxLength is exclusive)
2. Verify no overlapping tiers
3. Ensure tiers are in correct order (min to max)

### Issue: Discount not applied on top of boat length price

**Symptoms:**

- Boat tier price shows, but discount ignored
- `pricingMethod` is "BOAT_LENGTH" instead of "COMBINED"

**Solutions:**

1. Verify member has qualifying facilities rented
2. Check pricing rules are active
3. Ensure season parameter matches rented facilities

## Future Enhancements

1. **Database-Driven Tiers**
   - Store tiers in database table
   - Admin UI for tier management
   - Per-facility-type customization

2. **Dynamic Tier Adjustments**
   - Seasonal pricing variations
   - Market-based price updates
   - Inflation adjustments

3. **Additional Pricing Factors**
   - Boat width consideration
   - Engine type (motorboat vs sailboat)
   - Insurance requirements

4. **Analytics**
   - Track most common boat sizes
   - Optimize tier boundaries
   - Revenue analysis by tier

5. **Member Communication**
   - Show tier boundaries in UI
   - Price calculator widget
   - Email notifications of tier changes

## References

- `main/domain/facility_rental/pricing/boat_length_pricing.go`
- `main/domain/facility_rental/pricing/composite_price_calculator.go`
- `main/domain/facility_rental/rental_management_service.go`
- `main/infrastructure/http/handlers.go` (SuggestedPriceHandler)
- Original pricing system: `main/domain/facility_rental/pricing/suggested_price_calculator.go`

## Migration Guide

For existing installations:

1. **No database changes required** - Current implementation uses in-memory tiers
2. **API backward compatible** - `boat_length` parameter is optional
3. **Frontend updates optional** - System works without UI changes
4. **Existing prices preserved** - Base pricing still applies when boat_length not provided

To enable for existing facilities:

1. Ensure facility types have `HasBoat = true`
2. Restart backend to apply default tiers
3. Update frontend to send `boat_length` parameter
4. Test with various boat sizes

## Support

For questions or issues:

- Review this documentation
- Check troubleshooting section
- Examine unit tests for examples
- Contact backend development team
