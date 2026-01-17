# Migration Guide: Discount Percentages → Absolute Prices

This guide documents the changes made to the pricing system to switch from percentage-based discounts to absolute prices.

## What Changed?

The yacht club president requested that special prices be set as **absolute values** (e.g., "€80.00") rather than **percentage discounts** (e.g., "20% off"). This gives complete control over final pricing without dependency on base prices.

## Code Changes

### 1. Type Renaming

| Old Name | New Name | Reason |
|----------|----------|--------|
| `DiscountRule` | `PricingRule` | More accurate terminology |
| `FacilityTypeDiscountConfig` | `FacilityTypePricingConfig` | Reflects absolute pricing |
| `DiscountPercentage` | `SpecialPrice` | Now an absolute price value |
| `GetApplicableDiscounts()` | `GetApplicablePricingRules()` | Updated method name |
| `GetAllDiscountConfigs()` | `GetAllPricingConfigs()` | Updated method name |
| `SetDiscountConfigs()` | `SetPricingConfigs()` | Updated method name |

### 2. Struct Changes

**Before:**
```go
type DiscountRule struct {
    RequiredFacilityTypeId int64
    DiscountPercentage     float64  // e.g., 20.0 for 20% off
}
```

**After:**
```go
type PricingRule struct {
    RequiredFacilityTypeId int64
    SpecialPrice           float64  // e.g., 80.0 EUR
}
```

### 3. Calculation Logic

**Before (Percentage-based):**
```go
if rule.DiscountPercentage > bestDiscount {
    bestDiscount = rule.DiscountPercentage
}
// ...
discountAmount := baseSuggestedPrice * (bestDiscount / 100.0)
return baseSuggestedPrice - discountAmount
```

**After (Absolute pricing):**
```go
if !priceFound || rule.SpecialPrice < bestPrice {
    bestPrice = rule.SpecialPrice
    priceFound = true
}
// ...
return bestPrice  // Return the absolute special price
```

### 4. Configuration Changes

**Before:**
```go
{
    FacilityTypeId: 2,
    DiscountRules: []DiscountRule{
        {
            RequiredFacilityTypeId: 1,
            DiscountPercentage:     20.0,  // 20% off
        },
    },
}
```

**After:**
```go
{
    FacilityTypeId: 2,
    PricingRules: []PricingRule{
        {
            RequiredFacilityTypeId: 1,
            SpecialPrice:           80.0,  // €80.00
        },
    },
}
```

### 5. API Response Changes

**Before:**
```json
{
  "suggestedPrice": 80.00,
  "basePrice": 100.00,
  "discountApplied": 20.00,
  "discountPercentage": 20.0,
  "hasDiscount": true
}
```

**After:**
```json
{
  "suggestedPrice": 80.00,
  "basePrice": 100.00,
  "savingsAmount": 20.00,
  "hasSpecialPrice": true,
  "applicableRules": 1
}
```

### 6. Frontend Type Changes

**Before:**
```typescript
export interface SuggestedPriceResponse {
  suggestedPrice: number;
  basePrice: number;
  discountApplied: number;
  discountPercentage: number;
  hasDiscount: boolean;
}
```

**After:**
```typescript
export interface SuggestedPriceResponse {
  suggestedPrice: number;
  basePrice: number;
  savingsAmount: number;
  hasSpecialPrice: boolean;
  applicableRules: number;
}
```

## Migration Steps

### For Backend Code

1. **Update any code that references the old types:**
   - Replace `DiscountRule` → `PricingRule`
   - Replace `DiscountPercentage` → `SpecialPrice`
   - Replace `hasDiscount` → `hasSpecialPrice`

2. **Update pricing configurations:**
   - Convert percentage values to absolute prices
   - Example: If base price is €100 and discount was 20%, set `SpecialPrice: 80.0`

3. **Update method calls:**
   ```go
   // Old
   rules := service.GetApplicableDiscountsForMember(...)
   
   // New
   rules := service.GetApplicablePricingRulesForMember(...)
   ```

### For Frontend Code

1. **Update API response handling:**
   ```typescript
   // Old
   if (priceInfo.hasDiscount) {
       console.log(`${priceInfo.discountPercentage}% off!`);
   }
   
   // New
   if (priceInfo.hasSpecialPrice) {
       console.log(`Special price: €${priceInfo.suggestedPrice}`);
   }
   ```

2. **Update UI messages:**
   - "20% SCONTO" → "PREZZO SPECIALE"
   - "You get 20% off!" → "Special price available!"
   - Emphasize the final price rather than the percentage saved

3. **Update variable names:**
   ```typescript
   // Old
   let discountInfo = { hasDiscount: true, percentage: 20 };
   
   // New
   let priceInfo = { hasSpecialPrice: true, savingsAmount: 20.0 };
   ```

## Setting New Prices

To configure special prices for your facility types:

1. **Identify the facility type IDs** from your database
2. **Decide the special prices** (absolute values in EUR)
3. **Update the configuration** in `getDefaultPricingConfigs()`:

```go
return []FacilityTypePricingConfig{
    {
        FacilityTypeId: YOUR_FACILITY_TYPE_ID,
        PricingRules: []PricingRule{
            {
                RequiredFacilityTypeId: REQUIRED_FACILITY_ID,
                SpecialPrice:           YOUR_SPECIAL_PRICE,
            },
        },
    },
}
```

### Example Scenarios

**Scenario 1: Fixed price for members with a box**
- Base price for Canoe: €100
- Special price if you have a Box: €75
```go
{
    FacilityTypeId: 2, // Canoe
    PricingRules: []PricingRule{
        {RequiredFacilityTypeId: 1, SpecialPrice: 75.0},
    },
}
```

**Scenario 2: Different prices based on what you have**
- Base price for Locker: €70
- If you have a Box: €45
- If you have a Mooring: €55
- System picks the lowest price (€45 if you have both)
```go
{
    FacilityTypeId: 3, // Locker
    PricingRules: []PricingRule{
        {RequiredFacilityTypeId: 1, SpecialPrice: 45.0}, // Box
        {RequiredFacilityTypeId: 4, SpecialPrice: 55.0}, // Mooring
    },
}
```

## Benefits of This Change

1. **Full Control**: President sets exact prices, not influenced by base price changes
2. **Simplicity**: Easier to understand "€80" than "20% off €100"
3. **Flexibility**: Can set prices that don't follow a percentage pattern
4. **Clear Communication**: Members see the exact price they'll pay
5. **No Calculation Errors**: No rounding issues from percentage calculations

## Backwards Compatibility

⚠️ **This is a breaking change**. The old API response format is no longer supported.

If you have existing code that relies on the old format:
1. Update all API consumers to use the new response structure
2. Update any cached data or stored configurations
3. Test all pricing calculations thoroughly

## Testing Checklist

- [ ] Verify pricing calculator returns correct absolute prices
- [ ] Test with multiple applicable rules (lowest price wins)
- [ ] Test with no applicable rules (returns base price)
- [ ] Verify API response has correct structure
- [ ] Test frontend displays prices correctly
- [ ] Verify savings amount calculation
- [ ] Test with different facility type combinations

## Questions?

If you have questions about this migration, refer to:
- `README.md` for usage documentation
- `suggested_price_calculator.go` for implementation details
- `PRICING_EXAMPLE.md` for frontend integration examples
