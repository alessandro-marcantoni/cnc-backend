# Memberships API

## POST /api/v1.0/memberships

Add a new membership period for an existing member.

### Description

This endpoint creates a new membership period in the `membership_periods` table for an existing member. It associates the membership with a specific season and records the price paid.

### Request Body

```json
{
  "memberId": 123,
  "seasonId": 45,
  "seasonStartsAt": "2024-09-01T00:00:00Z",
  "seasonEndsAt": "2025-08-31T23:59:59Z",
  "price": 130.00
}
```

### Request Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `memberId` | `int64` | Yes | The ID of the member to add the membership for |
| `seasonId` | `int64` | Yes | The ID of the season for this membership period |
| `seasonStartsAt` | `string` | Yes | ISO 8601 timestamp for when the season starts |
| `seasonEndsAt` | `string` | Yes | ISO 8601 timestamp for when the season ends |
| `price` | `float64` | Yes | The price of the membership (must be > 0) |

### Response

**Status Code:** `201 Created`

Returns the complete member details including all memberships:

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

### Error Responses

#### 400 Bad Request

Invalid request body or missing required fields:

```json
{
  "error": "memberId is required"
}
```

```json
{
  "error": "price must be greater than 0"
}
```

#### 500 Internal Server Error

Database or server error:

```json
{
  "error": "failed to insert membership period: ..."
}
```

### Example Usage

```bash
curl -X POST http://localhost:8080/api/v1.0/memberships \
  -H "Content-Type: application/json" \
  -d '{
    "memberId": 123,
    "seasonId": 45,
    "seasonStartsAt": "2024-09-01T00:00:00Z",
    "seasonEndsAt": "2025-08-31T23:59:59Z",
    "price": 130.00
  }'
```

### Business Rules

1. The member must already exist in the system
2. The membership period is created with status `ACTIVE` (status_id = 1)
3. The membership_id is automatically retrieved from the existing membership record for the member
4. After successful creation, the endpoint returns the updated member details with all memberships

### Database Operations

1. Begins a database transaction
2. Retrieves the `membership_id` for the given `memberId`
3. Inserts a new record in `membership_periods` table
4. Commits the transaction
5. Fetches and returns the complete member details
