CREATE TYPE "order_status" AS ENUM (
  'PENDING',
  'PAYMENT_ACCEPTED',
  'ON_DELIVERY',
  'DELIVERED',
  'REJECTED',
  'FAILED',
  'CANCELLED'
);

CREATE TYPE "payment_status" AS ENUM (
  'PENDING',
  'APPROVED',
  'REJECTED'
);

CREATE TABLE IF NOT EXISTS "orders" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "status" order_status NOT NULL,
  "total_price" float NOT NULL,
  "payment_id" uuid,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "order_items" (
  "id" uuid PRIMARY KEY,
  "order_id" uuid,
  "product_id" uuid NOT NULL,
  "product_quantity" integer NOT NULL,
  "note" text,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "order_addresses" (
  "id" uuid PRIMARY KEY,
  "order_id" uuid,
  "street" varchar NOT NULL,
  "city" varchar NOT NULL,
  "state" varchar NOT NULL,
  "zip_code" varchar NOT NULL,
  "note" varchar,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "deleted_at" timestamp
);

CREATE INDEX ON "orders" ("status");
CREATE INDEX ON "order_items" ("order_id");
CREATE INDEX ON "order_addresses" ("order_id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
ALTER TABLE "order_addresses" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");