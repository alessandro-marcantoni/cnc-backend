# Frontend Integration Guide: Add Membership

## Overview
This guide shows how to integrate the `POST /api/v1.0/memberships` endpoint into the frontend application.

## API Endpoint

```
POST http://localhost:8080/api/v1.0/memberships
```

## TypeScript Interface

```typescript
interface AddMembershipRequest {
  memberId: number;
  seasonId: number;
  seasonStartsAt: string; // ISO 8601 format
  seasonEndsAt: string;   // ISO 8601 format
  price: number;
}

interface AddMembershipResponse {
  id: number;
  firstName: string;
  lastName: string;
  email: string;
  birthDate: string;
  phoneNumbers: PhoneNumber[];
  addresses: Address[];
  memberships: Membership[];
}
```

## Implementation Example

### 1. API Service Function

```typescript
// src/lib/api/memberships.ts
import { API_BASE_URL } from '$lib/config';

export async function addMembership(
  memberId: number,
  seasonId: number,
  seasonStartsAt: Date,
  seasonEndsAt: Date,
  price: number
): Promise<AddMembershipResponse> {
  const response = await fetch(`${API_BASE_URL}/api/v1.0/memberships`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      memberId,
      seasonId,
      seasonStartsAt: seasonStartsAt.toISOString(),
      seasonEndsAt: seasonEndsAt.toISOString(),
      price,
    }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to add membership');
  }

  return response.json();
}
```

### 2. Svelte Component Integration

```svelte
<script lang="ts">
  import { addMembership } from '$lib/api/memberships';
  import type { Season } from '$model/shared/season';
  
  interface Props {
    memberId: number;
    currentSeason: Season;
    onSuccess: () => void;
  }
  
  let { memberId, currentSeason, onSuccess }: Props = $props();
  
  let price = $state('130.00');
  let loading = $state(false);
  let error = $state<string | null>(null);
  
  async function handleAddMembership() {
    loading = true;
    error = null;
    
    try {
      await addMembership(
        memberId,
        currentSeason.id,
        new Date(currentSeason.startsAt),
        new Date(currentSeason.endsAt),
        parseFloat(price)
      );
      
      onSuccess();
    } catch (err) {
      error = err instanceof Error ? err.message : 'An error occurred';
    } finally {
      loading = false;
    }
  }
</script>

<button onclick={handleAddMembership} disabled={loading}>
  {loading ? 'Adding...' : 'Add Membership'}
</button>

{#if error}
  <p class="error">{error}</p>
{/if}
```

### 3. Integration with Renew Membership Dialog

Update the `renew-membership-dialog.svelte` to call the API:

```svelte
<script lang="ts">
  import { addMembership } from '$lib/api/memberships';
  
  // ... existing props
  
  async function onSubmit() {
    if (!selectedSeason || !price) return;
    
    try {
      const season = availableSeasons.find(s => s.name === selectedSeason);
      if (!season) throw new Error('Season not found');
      
      await addMembership(
        memberId,
        season.id,
        new Date(season.startsAt),
        new Date(season.endsAt),
        parseFloat(price)
      );
      
      // Success - close dialog and refresh data
      onClose();
      // Trigger parent component to refetch member data
      
    } catch (err) {
      console.error('Failed to add membership:', err);
      // Show error to user
    }
  }
</script>
```

## Request Payload Example

```json
{
  "memberId": 123,
  "seasonId": 45,
  "seasonStartsAt": "2024-09-01T00:00:00Z",
  "seasonEndsAt": "2025-08-31T23:59:59Z",
  "price": 130.00
}
```

## Response Example

```json
{
  "id": 123,
  "firstName": "Mario",
  "lastName": "Rossi",
  "email": "mario.rossi@example.com",
  "birthDate": "1980-05-15",
  "phoneNumbers": [
    {
      "prefix": "+39",
      "number": "3331234567"
    }
  ],
  "addresses": [
    {
      "country": "Italy",
      "city": "Rome",
      "street": "Via Roma",
      "streetNumber": "10",
      "zipCode": "00100"
    }
  ],
  "memberships": [
    {
      "id": 789,
      "number": 12345,
      "status": "ACTIVE",
      "validFrom": "2024-09-01",
      "expiresAt": "2025-08-31",
      "payment": null
    }
  ]
}
```

## Error Handling

### Common Errors

1. **400 Bad Request** - Missing or invalid fields
```json
{
  "error": "memberId is required"
}
```

2. **500 Internal Server Error** - Database or server error
```json
{
  "error": "failed to insert membership period: ..."
}
```

### Error Handling Example

```typescript
try {
  await addMembership(...);
} catch (error) {
  if (error instanceof Error) {
    if (error.message.includes('required')) {
      // Show validation error to user
      showToast('Please fill in all required fields', 'error');
    } else if (error.message.includes('failed to insert')) {
      // Show database error
      showToast('Unable to add membership. Please try again.', 'error');
    } else {
      // Generic error
      showToast('An unexpected error occurred', 'error');
    }
  }
}
```

## Complete Integration Flow

1. User clicks "Add Membership" button in `membership-card.svelte`
2. `renew-membership-dialog.svelte` opens
3. User selects season and enters price
4. User clicks "Confirm"
5. Frontend calls `POST /api/v1.0/memberships`
6. Backend creates membership period
7. Backend returns updated member details
8. Frontend closes dialog
9. Frontend refreshes member data to show new membership

## Data Refresh Pattern

After successfully adding a membership, refresh the member data:

```typescript
async function handleMembershipAdded() {
  try {
    // Add membership
    await addMembership(...);
    
    // Refetch member details to get updated data
    const updatedMember = await fetchMemberById(memberId, currentSeason.code);
    
    // Update local state
    memberDetails = updatedMember;
    
    // Close dialog
    dialogOpen = false;
    
    // Show success message
    showToast('Membership added successfully', 'success');
    
  } catch (error) {
    // Handle error
  }
}
```

## Testing with cURL

```bash
curl -X POST http://localhost:8080/api/v1.0/memberships \
  -H "Content-Type: application/json" \
  -d '{
    "memberId": 1,
    "seasonId": 1,
    "seasonStartsAt": "2024-09-01T00:00:00Z",
    "seasonEndsAt": "2025-08-31T23:59:59Z",
    "price": 130.00
  }'
```

## Notes

- Always convert dates to ISO 8601 format when sending to the backend
- The `memberId` must reference an existing member
- The `seasonId` must reference an existing season in the database
- Price must be greater than 0
- The backend automatically sets the membership status to ACTIVE
- After adding, the member will have access to the new membership period
