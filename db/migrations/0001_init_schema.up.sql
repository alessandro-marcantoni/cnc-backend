-- =========================
-- MEMBERS
-- =========================
CREATE TABLE IF NOT EXISTS members (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name  VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE
);

-- =========================
-- PHONE NUMBERS (1:N)
-- =========================
CREATE TABLE IF NOT EXISTS phone_numbers (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    number VARCHAR(50) NOT NULL,
    description VARCHAR(255),
    UNIQUE(member_id, number)
);

CREATE INDEX IF NOT EXISTS idx_phone_numbers_member
ON phone_numbers(member_id);

-- =========================
-- ADDRESSES (1:N)
-- =========================
CREATE TABLE IF NOT EXISTS addresses (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    country VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    street VARCHAR(255) NOT NULL,
    street_number VARCHAR(50) NOT NULL,
    zip_code VARCHAR(20) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_addresses_member
ON addresses(member_id);

-- =========================
-- FACILITIES CATALOG
-- =========================
CREATE TABLE IF NOT EXISTS facilities_catalog (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    suggested_price NUMERIC(10,2) NOT NULL CHECK (suggested_price >= 0)
);

CREATE TABLE IF NOT EXISTS facilities (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    facility_type_id BIGINT NOT NULL REFERENCES facilities_catalog(id),
    identifier VARCHAR(255) NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_facilities_type
ON facilities(facility_type_id);


-- =========================
-- SEASONS
-- =========================
CREATE TABLE IF NOT EXISTS seasons (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,          -- e.g. '2025', '2025-2026'
    name VARCHAR(100) NOT NULL,                -- e.g. 'Season 2025'
    starts_at DATE NOT NULL,
    ends_at   DATE NOT NULL,
    CHECK (ends_at > starts_at)
);

-- =========================
-- RENTED FACILITIES (HISTORY)
-- =========================
CREATE TABLE IF NOT EXISTS rented_facilities (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    facility_id BIGINT NOT NULL REFERENCES facilities(id),
    member_id BIGINT NOT NULL REFERENCES members(id),
    season_id BIGINT NOT NULL REFERENCES seasons(id),
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'EUR',
    deleted_at TIMESTAMP DEFAULT NULL
);

-- Partial unique index: only applies to non-deleted records
-- This allows re-renting the same facility in the same season after soft delete
CREATE UNIQUE INDEX IF NOT EXISTS idx_rented_facilities_active_unique
ON rented_facilities(facility_id, season_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_rented_facility_member
ON rented_facilities(member_id);

CREATE INDEX IF NOT EXISTS idx_rented_facility_facility
ON rented_facilities(facility_id);

CREATE INDEX IF NOT EXISTS idx_rented_facilities_deleted_at
ON rented_facilities(deleted_at) WHERE deleted_at IS NULL;

-- =========================
-- MEMBERSHIPS
-- =========================
CREATE TABLE IF NOT EXISTS membership_statuses (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    status VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS memberships (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    member_id BIGINT NOT NULL REFERENCES members(id),
    number BIGINT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_memberships_member
ON memberships(member_id);

CREATE TABLE IF NOT EXISTS membership_periods (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    membership_id BIGINT NOT NULL
        REFERENCES memberships(id) ON DELETE CASCADE,
    status_id BIGINT NOT NULL
        REFERENCES membership_statuses(id),
    exclusion_deliberated_at TIMESTAMP,
    season_id BIGINT NOT NULL REFERENCES seasons(id),
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'EUR',
    UNIQUE (membership_id, season_id)
);

CREATE INDEX IF NOT EXISTS idx_membership_periods_membership
ON membership_periods(membership_id);

CREATE INDEX IF NOT EXISTS idx_membership_periods_status
ON membership_periods(status_id);

-- =========================
-- PAYMENTS (POLYMORPHIC)
-- =========================
CREATE TABLE IF NOT EXISTS payments (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    rented_facility_id BIGINT REFERENCES rented_facilities(id),
    membership_period_id BIGINT REFERENCES membership_periods(id),
    amount NUMERIC(10,2) NOT NULL CHECK (amount >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'EUR',
    paid_at TIMESTAMP NOT NULL DEFAULT now(),
    payment_method VARCHAR(255) NOT NULL,
    transaction_ref TEXT UNIQUE,
    CHECK (
        (rented_facility_id IS NOT NULL AND membership_period_id IS NULL)
     OR (rented_facility_id IS NULL AND membership_period_id IS NOT NULL)
    )
);

CREATE INDEX IF NOT EXISTS idx_payments_membership_period
ON payments(membership_period_id);

CREATE INDEX IF NOT EXISTS idx_payments_rented_facility
ON payments(rented_facility_id);

CREATE INDEX IF NOT EXISTS idx_payments_paid_at
ON payments(paid_at);

-- =========================
-- BOATS (1:1 WITH RENTED)
-- =========================
CREATE TABLE IF NOT EXISTS boats (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    rented_facility_id BIGINT NOT NULL UNIQUE
        REFERENCES rented_facilities(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    length_meters NUMERIC(10,2) NOT NULL CHECK (length_meters > 0),
    width_meters  NUMERIC(10,2) NOT NULL CHECK (width_meters > 0)
);

-- =========================
-- INSURANCES (1:N WITH BOAT)
-- =========================
CREATE TABLE IF NOT EXISTS insurances (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    boat_id BIGINT NOT NULL REFERENCES boats(id) ON DELETE CASCADE,
    provider VARCHAR(255) NOT NULL,
    number VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_insurances_boat
ON insurances(boat_id);

CREATE INDEX IF NOT EXISTS idx_insurances_expiry
ON insurances(expires_at);

-- =========================
-- WAITING LIST
-- =========================
CREATE TABLE IF NOT EXISTS members_waiting (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    member_id BIGINT NOT NULL REFERENCES members(id),
    facility_type_id BIGINT NOT NULL REFERENCES facilities_catalog(id),
    queued_at TIMESTAMP NOT NULL DEFAULT now(),
    notes TEXT,
    UNIQUE(member_id, facility_type_id)
);

CREATE INDEX IF NOT EXISTS idx_waiting_member
ON members_waiting(member_id);

CREATE INDEX IF NOT EXISTS idx_waiting_facility_type
ON members_waiting(facility_type_id);

CREATE INDEX IF NOT EXISTS idx_waiting_queued
ON members_waiting(queued_at);

-- =========================
-- FACILITY PRICING RULES
-- =========================
-- This table stores special pricing rules that apply when a member
-- already has a specific facility type rented.
-- Example: If member has a Box, then Canoe costs 80 EUR instead of base price
CREATE TABLE IF NOT EXISTS facility_pricing_rules (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    facility_type_id BIGINT NOT NULL REFERENCES facilities_catalog(id) ON DELETE CASCADE,
    required_facility_type_id BIGINT NOT NULL REFERENCES facilities_catalog(id) ON DELETE CASCADE,
    special_price NUMERIC(10,2) NOT NULL CHECK (special_price >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'EUR',
    description TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    -- Ensure we don't have duplicate rules for the same combination
    UNIQUE(facility_type_id, required_facility_type_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_pricing_rules_facility_type
ON facility_pricing_rules(facility_type_id);

CREATE INDEX IF NOT EXISTS idx_pricing_rules_required_facility_type
ON facility_pricing_rules(required_facility_type_id);

CREATE INDEX IF NOT EXISTS idx_pricing_rules_active
ON facility_pricing_rules(active);

-- =========================
-- BETTER AUTH TABLES
-- =========================
create table "user" ("id" text not null primary key, "name" text not null, "email" text not null unique, "emailVerified" boolean not null, "image" text, "createdAt" timestamptz default CURRENT_TIMESTAMP not null, "updatedAt" timestamptz default CURRENT_TIMESTAMP not null);

create table "session" ("id" text not null primary key, "expiresAt" timestamptz not null, "token" text not null unique, "createdAt" timestamptz default CURRENT_TIMESTAMP not null, "updatedAt" timestamptz not null, "ipAddress" text, "userAgent" text, "userId" text not null references "user" ("id") on delete cascade);

create table "account" ("id" text not null primary key, "accountId" text not null, "providerId" text not null, "userId" text not null references "user" ("id") on delete cascade, "accessToken" text, "refreshToken" text, "idToken" text, "accessTokenExpiresAt" timestamptz, "refreshTokenExpiresAt" timestamptz, "scope" text, "password" text, "createdAt" timestamptz default CURRENT_TIMESTAMP not null, "updatedAt" timestamptz not null);

create table "verification" ("id" text not null primary key, "identifier" text not null, "value" text not null, "expiresAt" timestamptz not null, "createdAt" timestamptz default CURRENT_TIMESTAMP not null, "updatedAt" timestamptz default CURRENT_TIMESTAMP not null);

create index "session_userId_idx" on "session" ("userId");

create index "account_userId_idx" on "account" ("userId");

create index "verification_identifier_idx" on "verification" ("identifier");

create table "jwks" ("id" text not null primary key, "publicKey" text not null, "privateKey" text not null, "createdAt" timestamptz not null, "expiresAt" timestamptz);
