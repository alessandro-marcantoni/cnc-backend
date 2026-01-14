# Quick Reference: Create Member API

## Endpoint

```
POST /api/v1.0/members
```

## Minimal Request (No Membership)

```json
{
  "firstName": "Anna",
  "lastName": "Bianchi",
  "birthDate": "1990-07-22",
  "email": "anna.bianchi@example.com",
  "phoneNumbers": [],
  "addresses": [],
  "createMembership": false
}
```

## Full Request (With Membership)

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

## Response (201 Created)

```json
{
  "id": 123,
  "firstName": "Mario",
  "lastName": "Rossi",
  "email": "mario.rossi@example.com",
  "birthDate": "1985-03-15",
  "phoneNumbers": [...],
  "addresses": [...],
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

## cURL Example

```bash
curl -X POST http://localhost:8080/api/v1.0/members \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Mario",
    "lastName": "Rossi",
    "birthDate": "1985-03-15",
    "email": "mario.rossi@example.com",
    "phoneNumbers": [{"prefix": "+39", "number": "3331234567"}],
    "addresses": [{"country": "Italy", "city": "Rome", "street": "Via Roma", "streetNumber": "123", "zipCode": "00100"}],
    "createMembership": true,
    "seasonId": 1,
    "price": 130.0
  }'
```

## Field Constraints

| Field                    | Type    | Required | Validation                        |
| ------------------------ | ------- | -------- | --------------------------------- |
| firstName                | string  | ✓        | -                                 |
| lastName                 | string  | ✓        | -                                 |
| birthDate                | string  | ✓        | Format: YYYY-MM-DD                |
| email                    | string  | ✓        | Valid email, unique               |
| phoneNumbers[].prefix    | string  | -        | 3 characters if provided          |
| phoneNumbers[].number    | string  | ✓        | 10-15 digits                      |
| addresses[].country      | string  | ✓        | -                                 |
| addresses[].city         | string  | ✓        | -                                 |
| addresses[].street       | string  | ✓        | -                                 |
| addresses[].streetNumber | string  | ✓        | -                                 |
| addresses[].zipCode      | string  | -        | -                                 |
| createMembership         | boolean | ✓        | -                                 |
| seasonId                 | integer | -        | Required if createMembership=true |
| price                    | number  | -        | Custom price, defaults to €130.00 |

## What Gets Created

### When `createMembership: false`

- ✓ Member record
- ✓ Phone numbers
- ✓ Addresses

### When `createMembership: true`

- ✓ Member record
- ✓ Phone numbers
- ✓ Addresses
- ✓ Membership with auto-generated number
- ✓ Membership period (ACTIVE, custom or default price, season dates)

## Common Errors

| Status | Error                | Solution                                    |
| ------ | -------------------- | ------------------------------------------- |
| 400    | invalid email format | Use valid email address                     |
| 400    | invalid birth date   | Use YYYY-MM-DD format                       |
| 400    | invalid phone number | 10-15 digits required                       |
| 400    | seasonId required    | Provide seasonId when createMembership=true |
| 500    | duplicate email      | Email already exists                        |
| 500    | invalid season       | Check seasonId exists                       |

## Transaction Guarantee

✓ All operations are atomic  
✓ Failure = complete rollback  
✓ Success = all data saved

## Default Values

- Membership Status: **ACTIVE**
- Membership Price: **€130.00** (unless custom price provided)
- Currency: **EUR**
- Validity Period: **Season dates**

## Quick Test

```bash
# Test without membership
curl -X POST http://localhost:8080/api/v1.0/members \
  -H "Content-Type: application/json" \
  -d '{"firstName":"Test","lastName":"User","birthDate":"1990-01-01","email":"test@example.com","phoneNumbers":[],"addresses":[],"createMembership":false}'

# Test with membership (default price)
curl -X POST http://localhost:8080/api/v1.0/members \
  -H "Content-Type: application/json" \
  -d '{"firstName":"Test","lastName":"User","birthDate":"1990-01-01","email":"test2@example.com","phoneNumbers":[],"addresses":[],"createMembership":true,"seasonId":1}'

# Test with membership (custom price)
curl -X POST http://localhost:8080/api/v1.0/members \
  -H "Content-Type: application/json" \
  -d '{"firstName":"Test","lastName":"User","birthDate":"1990-01-01","email":"test3@example.com","phoneNumbers":[],"addresses":[],"createMembership":true,"seasonId":1,"price":150.0}'
```
