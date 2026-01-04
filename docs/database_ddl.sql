CREATE TABLE "members" (
  "id" BIGINT NOT NULL,
  "first_name" VARCHAR(255) NOT NULL,
  "last_name" VARCHAR(255) NOT NULL,
  "date_of_birth" DATE NOT NULL,
  "email" VARCHAR(255) NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE "phone_numbers" (
  "id" BIGINT NOT NULL,
  "member_id" BIGINT NOT NULL,
  "number" VARCHAR(255) NOT NULL,
  "description" VARCHAR(255) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "FK_phone_numbers_member_id"
    FOREIGN KEY ("member_id")
      REFERENCES "members"("id")
);

CREATE TABLE "addresses" (
  "id" BIGINT NOT NULL,
  "member_id" BIGINT NOT NULL,
  "country" VARCHAR(255) NOT NULL,
  "city" VARCHAR(255) NOT NULL,
  "street" VARCHAR(255) NOT NULL,
  "number" VARCHAR(255) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "FK_addresses_member_id"
    FOREIGN KEY ("member_id")
      REFERENCES "members"("id")
);

CREATE TABLE "service_types" (
  "id" BIGINT NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  "description" VARCHAR(255) NOT NULL,
  "suggested_price" DECIMAL(10,2) NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE "services" (
  "id" BIGINT NOT NULL,
  "service_type" BIGINT NOT NULL,
  "identifier" VARCHAR(255) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "FK_services_service_type"
    FOREIGN KEY ("service_type")
      REFERENCES "service_types"("id")
);

CREATE TABLE "rented_services" (
  "id" BIGINT NOT NULL,
  "service_id" BIGINT NOT NULL,
  "member_id" BIGINT NOT NULL,
  "expires_at" DATETIME NOT NULL,
  "rented_at" DATETIME NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "FK_rented_services_member_id"
    FOREIGN KEY ("member_id")
      REFERENCES "members"("id"),
  CONSTRAINT "FK_rented_services_service_id"
    FOREIGN KEY ("service_id")
      REFERENCES "services"("id")
);

CREATE TABLE "membership_statuses" (
  "id" BIGINT NOT NULL,
  "status" VARCHAR(255) NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE "memberships" (
  "id" BIGINT NOT NULL,
  "member_id" BIGINT NOT NULL,
  "status" BIGINT NOT NULL,
  "number" BIGINT NOT NULL,
  "expires_at" DATETIME NOT NULL,
  "exclusion_deliberated_at" DATETIME NOT NULL,
  "excluded_at" DATETIME NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "FK_memberships_status"
    FOREIGN KEY ("status")
      REFERENCES "membership_statuses"("id"),
  CONSTRAINT "FK_memberships_member_id"
    FOREIGN KEY ("member_id")
      REFERENCES "members"("id")
);

CREATE TABLE "payments" (
  "id" BIGINT NOT NULL,
  "rented_service_id" BIGINT NOT NULL,
  "membership_id" BIGINT NOT NULL,
  "amount" DECIMAL(10,2) NOT NULL,
  "paid" BOOLEAN NOT NULL,
  "paid_at" DATETIME NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "FK_payments_membership_id"
    FOREIGN KEY ("membership_id")
      REFERENCES "memberships"("id"),
  CONSTRAINT "FK_payments_rented_service_id"
    FOREIGN KEY ("rented_service_id")
      REFERENCES "rented_services"("id"),
  CHECK (
          (rented_service_id IS NOT NULL AND membership_id IS NULL)
        OR (rented_service_id IS NULL AND membership_id IS NOT NULL)
    )
);

CREATE TABLE "boats" (
  "id" BIGINT NOT NULL,
  "rented_service_id" BIGINT NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  "length_meters" DECIMAL(10,2) NOT NULL,
  "width_meters" DECIMAL(10,2) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "FK_boats_rented_service_id"
    FOREIGN KEY ("rented_service_id")
      REFERENCES "rented_services"("id")
);

CREATE TABLE "insurances" (
  "id" BIGINT NOT NULL,
  "boat_id" BIGINT NOT NULL,
  "provider" VARCHAR(255) NOT NULL,
  "number" VARCHAR(255) NOT NULL,
  "expires_at" DATETIME NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "FK_insurances_boat_id"
    FOREIGN KEY ("boat_id")
      REFERENCES "boats"("id")
);

CREATE TABLE "members_waiting" (
  "id" BIGINT NOT NULL,
  "member_id" BIGINT NOT NULL,
  "service_type_id" BIGINT NOT NULL,
  "queued_at" DATETIME NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "FK_members_waiting_member_id"
    FOREIGN KEY ("member_id")
      REFERENCES "members"("id"),
  CONSTRAINT "FK_members_waiting_service_type"
    FOREIGN KEY ("service_type")
      REFERENCES "service_types"("id")
);
