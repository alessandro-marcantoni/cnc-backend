# Create Member API Endpoint

## Overview

This endpoint creates a new member in the system with their personal information, contact details, and optionally creates a membership with a membership period for the specified season.

## Endpoint

```
POST /api/v1.0/members
```

## Request Body

```json
{
  "firstName": "Mario",
  "lastName": "Rossi",
  "birthDate": "1985-03-15",
  "email": "mario.rossi@example.com",
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
      "streetNumber": "123",
      "zipCode": "00100"
    }
  ],
  "createMembership": true,
  "seasonId": 1,
  "price": 130.0
}
```

## Request Fields

| Field                      | Type    | Required | Description                                                                    |
| -------------------------- | ------- | -------- | ------------------------------------------------------------------------------ |
| `firstName`                | string  | Yes      | Member's first name                                                            |
| `lastName`                 | string  | Yes      | Member's last name                                                             |
| `birthDate`                | string  | Yes      | Member's date of birth in ISO format (YYYY-MM-DD)                              |
| `email`                    | string  | Yes      | Member's email address (must be valid and unique)                              |
| `phoneNumbers`             | array   | No       | Array of phone numbers                                                         |
| `phoneNumbers[].prefix`    | string  | No       | International prefix (e.g., "+39")                                             |
| `phoneNumbers[].number`    | string  | Yes      | Phone number (10-15 digits)                                                    |
| `addresses`                | array   | No       | Array of addresses                                                             |
| `addresses[].country`      | string  | Yes      | Country name                                                                   |
| `addresses[].city`         | string  | Yes      | City name                                                                      |
| `addresses[].street`       | string  | Yes      | Street name                                                                    |
| `addresses[].streetNumber` | string  | Yes      | Street number                                                                  |
| `addresses[].zipCode`      | string  | No       | Postal/ZIP code                                                                |
| `createMembership`         | boolean | Yes      | Whether to create a membership for this member                                 |
| `seasonId`                 | integer | No\*     | Season ID for the membership period (\*required if `createMembership` is true) |
| `price`                    | number  | No       | Custom membership price (defaults to €130.00 if not provided)                  |

## Response

### Success Response (201 Created)

```json
{
  "id": 123,
  "firstName": "Mario",
  "lastName": "Rossi",
  "email": "mario.rossi@example.com",
  "birthDate": "1985-03-15",
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
      "streetNumber": "123",
      "zipCode": "00100"
    }
  ],
  "memberships": [
    {
      "id": 456,
      "number": 1001,
      "status": "ACTIVE",
      "validFrom": "2025-04-01",
      "expiresAt": "2026-03-31",
      "payment": null
    }
  ]
}
```

### Error Responses

#### 400 Bad Request

```json
{
  "error": "invalid email: invalid email format"
}
```

Possible error messages:

- `"invalid JSON: <details>"` - Malformed JSON in request body
- `"invalid birth date: <details>"` - Birth date not in YYYY-MM-DD format
- `"invalid email: <details>"` - Email validation failed
- `"invalid phone number: <details>"` - Phone number validation failed
- `"seasonId is required when createMembership is true"` - Missing seasonId when trying to create membership

#### 500 Internal Server Error

```json
{
  "error": "failed to insert member: <details>"
}
```

Possible error messages:

- `"service not initialized"` - Server configuration error
- `"failed to begin transaction: <details>"` - Database transaction error
- `"failed to insert member: <details>"` - Member insertion failed (e.g., duplicate email)
- `"failed to insert phone number: <details>"` - Phone number insertion failed
- `"failed to insert address: <details>"` - Address insertion failed
- `"failed to get next membership number: <details>"` - Membership number generation failed
- `"failed to insert membership: <details>"` - Membership insertion failed
- `"failed to get season dates: <details>"` - Invalid season ID
- `"failed to insert membership period: <details>"` - Membership period insertion failed
- `"failed to commit transaction: <details>"` - Transaction commit failed

## Business Logic

### Transaction Flow

The endpoint executes the following operations within a database transaction:

1. **Insert Member** - Creates a new record in the `members` table with personal information
2. **Insert Phone Numbers** - Creates records in the `phone_numbers` table for each provided phone number
3. **Insert Addresses** - Creates records in the `addresses` table for each provided address
4. **Create Membership** (if `createMembership` is true):
   - Get the next available membership number (max current number + 1)
   - Insert a new record in the `memberships` table
   - Get the season start and end dates
   - Insert a membership period in the `membership_periods` table with:
     - Status: ACTIVE (status_id = 1)
     - Validity: Season start date to season end date
     - Price: Custom price if provided, otherwise default suggested membership price (€130.00)
5. **Commit Transaction** - All operations succeed or all are rolled back
6. **Fetch Member Details** - Returns the complete member details including all memberships

### Validation Rules

- **Email**: Must be a valid email format and unique in the system
- **Phone Number**:
  - Number must be 10-15 characters
  - Prefix must be exactly 3 characters (if provided)
- **Birth Date**: Must be in ISO format (YYYY-MM-DD)
- **Season**: Must exist in the database
- **Membership Number**: Automatically generated as sequential number

### Default Values

- **Membership Status**: ACTIVE when a membership is created
- **Membership Price**: €130.00 (defined in `membership.SuggestedMembershipPrice`) unless custom price is provided
- **Membership Validity**: Matches the season dates (starts_at to ends_at)
- **Currency**: EUR (Euro)

## Database Schema Impact

This endpoint affects the following tables:

- `members` - One new row
- `phone_numbers` - One row per phone number provided
- `addresses` - One row per address provided
- `memberships` - One row if `createMembership` is true
- `membership_periods` - One row if `createMembership` is true

## Examples

### Example 1: Create member without membership

```bash
curl -X POST http://localhost:8080/api/v1.0/members \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Anna",
    "lastName": "Bianchi",
    "birthDate": "1990-07-22",
    "email": "anna.bianchi@example.com",
    "phoneNumbers": [],
    "addresses": [],
    "createMembership": false
  }'
```

### Example 2: Create member with membership for Season 2025

```bash
curl -X POST http://localhost:8080/api/v1.0/members \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Giovanni",
    "lastName": "Verdi",
    "birthDate": "1982-11-30",
    "email": "giovanni.verdi@example.com",
    "phoneNumbers": [
      {
        "prefix": "+39",
        "number": "3487654321"
      },
      {
        "prefix": "",
        "number": "0612345678"
      }
    ],
    "addresses": [
      {
        "country": "Italy",
        "city": "Milan",
        "street": "Corso Buenos Aires",
        "streetNumber": "45",
        "zipCode": "20124"
      }
    ],
    "createMembership": true,
    "seasonId": 1,
    "price": 130.00
  }'
```

## Notes

- All operations are atomic - if any step fails, the entire operation is rolled back
- Email addresses must be unique across all members
- Membership numbers are automatically generated and sequential
- The created member is immediately returned with all their details
- If creating a membership, the member will have an ACTIVE status by default
- The membership period will not include a payment initially (payment can be added separately)
- `seasonId` is optional but required when `createMembership` is `true`
- `price` is optional and defaults to €130.00 if not provided
- If `createMembership` is `false`, member is created without any membership data
