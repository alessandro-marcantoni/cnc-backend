# Boat Length-Based Pricing - Quick Start Guide

## What Is This?

A new pricing method that calculates facility rental prices based on the length of the boat being stored. Larger boats = higher prices.

## Quick Example

```bash
# Without boat length - gets base price
curl "http://localhost:8080/api/v1.0/suggested-price?facility_type_id=2&member_id=5&season=1"
# Response: {"suggestedPrice": 100.0, "pricingMethod": "BASE"}

# With boat length - gets tier price
curl "http://localhost:8080/api/v1.0/suggested-price?facility_type_id=2&member_id=5&season=1&boat_length=7.5"
# Response: {"suggestedPrice": 130.0, "pricingMethod": "BOAT_LENGTH"}
```

## Default Pricing Tiers

| Boat Length | Price Multiplier | Example (Base: €100) |
|-------------|------------------|---------------------|
| 0m - 6m     | 1.0x            | €100                |
| 6m - 8m     | 1.3x            | €130                |
| 8m - 10m    | 1.6x            | €160                |
| 10m+        | 2.0x            | €200                |

## How to Use (Backend)

### Option 1: Service Layer

```go
import (
    "github.com/alessandro-marcantoni/cnc-backend/main/domain"
    "github.com/alessandro-marcantoni/cnc-backend/main/domain/facility_rental"
)

// Calculate price with boat length
boatLength := 7.5
result := rentalService.GetSuggestedPriceWithBoatLength(
    facilityTypeId,
    100.0, // base price
    memberId,
    seasonId,
    &boatLength,
)

fmt.Printf("Final Price: €%.2f\n", result.FinalPrice)
fmt.Printf("Method: %s\n", result.PricingMethod)
```

### Option 2: Direct Calculator

```go
import "github.com/alessandro-marcantoni/cnc-backend/main/domain/facility_rental/pricing"

// Create context
ctx := pricing.PriceCalculationContext{
    FacilityTypeId:            2,
    BaseSuggestedPrice:        100.0,
    MemberRentedFacilityTypes: []int64{1, 3},
    BoatLengthMeters:          &boatLength,
}

// Calculate
result := compositePriceCalculator.CalculatePrice(ctx)
```

## How to Use (Frontend)

### JavaScript/TypeScript

```typescript
// Fetch pricing with boat length
async function getPrice(facilityTypeId: number, memberId: number, boatLength: number) {
  const url = new URL('/api/v1.0/suggested-price', window.location.origin);
  url.searchParams.set('facility_type_id', facilityTypeId.toString());
  url.searchParams.set('member_id', memberId.toString());
  url.searchParams.set('season', seasonId.toString());
  url.searchParams.set('boat_length', boatLength.toString());

  const response = await fetch(url);
  return await response.json();
}

// Use it
const pricing = await getPrice(2, 5, 7.5);
console.log(`Price: €${pricing.suggestedPrice}`);
console.log(`Method: ${pricing.pricingMethod}`);
```

### Svelte

```svelte
<script>
  let boatLength = 7.5;
  let pricing = $state({ suggestedPrice: 0 });

  async function updatePrice() {
    const response = await fetch(
      `/api/v1.0/suggested-price?facility_type_id=2&member_id=5&boat_length=${boatLength}`
    );
    pricing = await response.json();
  }

  $effect(() => {
    boatLength;
    updatePrice();
  });
</script>

<input type="number" bind:value={boatLength} step="0.1" min="0" />
<p>Price: €{pricing.suggestedPrice.toFixed(2)}</p>
```

## API Response Structure

```json
{
  "suggestedPrice": 130.0,          // Final calculated price
  "basePrice": 100.0,                // Original base price
  "pricingMethod": "BOAT_LENGTH",    // BASE | DISCOUNT | BOAT_LENGTH | COMBINED
  "discountApplied": false,          // Was discount applied?
  "discountAmount": 0,               // Amount saved from discount
  "boatLengthTierApplied": true,     // Was boat tier applied?
  "boatLengthTierPrice": 130.0,      // Price from boat tier
  "applicableRules": 0,              // Number of discount rules
  "boatLengthTiers": [               // All available tiers
    {
      "minLengthMeters": 0,
      "maxLengthMeters": 6.0,
      "price": 100.0
    },
    {
      "minLengthMeters": 6.0,
      "maxLengthMeters": 8.0,
      "price": 130.0
    }
  ],
  "hasBoatLengthPricing": true // Is boat pricing available?
}
```

## Pricing Priority

The system applies strategies in this order:

1. **Boat Length** (if boat facility + length provided) → Base tier price
2. **Discount** (if member qualifies) → Applied on top of tier price
3. **Base Price** (fallback) → Original suggested price

### Example: Combined Pricing

Member has a 7m boat AND owns a Box (which gives 20% discount):

1. Boat length 7m → Tier 6-8m → €130
2. Owns Box → 20% discount → -€26
3. **Final: €104**

## Test It

### 1. Start Backend

```bash
cd cnc-backend
go run .
```

### 2. Test Small Boat (5m)

```bash
curl "http://localhost:8080/api/v1.0/suggested-price?facility_type_id=2&member_id=1&boat_length=5"
```

Expected: €100 (first tier)

### 3. Test Medium Boat (7m)

```bash
curl "http://localhost:8080/api/v1.0/suggested-price?facility_type_id=2&member_id=1&boat_length=7"
```

Expected: €130 (second tier)

### 4. Test Large Boat (15m)

```bash
curl "http://localhost:8080/api/v1.0/suggested-price?facility_type_id=2&member_id=1&boat_length=15"
```

Expected: €200 (unlimited tier)

## Common Issues

### ❌ Price doesn't change with boat length

**Problem:** Facility type doesn't have `HasBoat = true`

**Solution:** Check facility type configuration:

```sql
SELECT id, name, has_boat FROM facilities_catalog WHERE id = 2;
```

### ❌ Getting base price instead of tier price

**Problem:** `boat_length` parameter missing or zero

**Solution:** Ensure you're sending the parameter:

```javascript
// ❌ Wrong
fetch('/api/v1.0/suggested-price?facility_type_id=2&member_id=1')

// ✅ Correct
fetch('/api/v1.0/suggested-price?facility_type_id=2&member_id=1&boat_length=7.5')
```

### ❌ Discount not combining with boat length price

**Problem:** Member doesn't have qualifying facility for the season

**Solution:** Verify member's rented facilities:

```bash
curl "http://localhost:8080/api/v1.0/rented-facilities?member_id=1&season=1"
```

## Customization

### Change Tier Prices

Edit `rental_management_service.go` in `buildBoatLengthPricingConfigs()`:

```go
tiers := []pricing.BoatLengthTier{
    {
        MinLengthMeters: 0,
        MaxLengthMeters: 6.0,
        Price:           basePrice * 1.0,  // Change multiplier here
    },
    {
        MinLengthMeters: 6.0,
        MaxLengthMeters: 8.0,
        Price:           basePrice * 1.5,  // Increase to 1.5x
    },
    // Add more tiers...
}
```

### Add New Tier

```go
{
    MinLengthMeters: 8.0,
    MaxLengthMeters: 12.0,
    Price:           basePrice * 1.8,
},
```

## Frontend Integration Example

### Complete Rent Facility Form

```svelte
<script>
  let facilityTypeId = 2;
  let memberId = 5;
  let boatLength = 7.5;
  let pricing = $state(null);
  let loading = $state(false);

  async function updatePrice() {
    loading = true;
    try {
      const response = await fetch(
        `/api/v1.0/suggested-price?` +
        `facility_type_id=${facilityTypeId}&` +
        `member_id=${memberId}&` +
        `boat_length=${boatLength}`
      );
      pricing = await response.json();
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    boatLength;
    updatePrice();
  });
</script>

<div class="form">
  <h2>Rent Boat Facility</h2>

  <label>
    Boat Length (meters):
    <input type="number" bind:value={boatLength} step="0.1" min="0" max="50" />
  </label>

  {#if loading}
    <p>Calculating price...</p>
  {:else if pricing}
    <div class="pricing-info">
      <h3>Price: €{pricing.suggestedPrice.toFixed(2)}</h3>

      {#if pricing.boatLengthTierApplied}
        <p>✓ Boat length tier: €{pricing.boatLengthTierPrice.toFixed(2)}</p>
      {/if}

      {#if pricing.discountApplied}
        <p>✓ Discount: -€{pricing.discountAmount.toFixed(2)}</p>
      {/if}

      {#if pricing.hasBoatLengthPricing}
        <details>
          <summary>View all pricing tiers</summary>
          <ul>
            {#each pricing.boatLengthTiers as tier}
              <li class:active={boatLength >= tier.minLengthMeters && boatLength < (tier.maxLengthMeters ?? 999)}>
                {tier.minLengthMeters}m - {tier.maxLengthMeters ?? '∞'}m: €{tier.price.toFixed(2)}
              </li>
            {/each}
          </ul>
        </details>
      {/if}
    </div>
  {/if}

  <button onclick={confirmRental}>Confirm Rental</button>
</div>

<style>
  .active {
    font-weight: bold;
    color: green;
  }
</style>
```

## Run Tests

```bash
# Run all pricing tests
cd cnc-backend
go test ./test/domain/facility_rental/pricing/... -v

# Run specific test
go test ./test/domain/facility_rental/pricing/... -run TestBoatLengthPriceCalculator -v
```

## Next Steps

1. **Try it out** - Use the curl examples above
2. **Integrate frontend** - Add boat length input to rent facility form
3. **Customize tiers** - Adjust multipliers to match your pricing strategy
4. **Add validation** - Ensure boat length is reasonable (0-50m)
5. **Show tiers to users** - Display pricing table in UI

## Resources

- Full Documentation: `docs/BOAT_LENGTH_PRICING.md`
- Implementation Summary: `docs/BOAT_LENGTH_PRICING_SUMMARY.md`
- Unit Tests: `test/domain/facility_rental/pricing/boat_length_pricing_test.go`
- Code: `main/domain/facility_rental/pricing/`

## Need Help?

- Check the full documentation for detailed examples
- Review unit tests for usage patterns
- Contact backend development team
