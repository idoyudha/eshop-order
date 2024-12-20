CREATE TYPE "order_status" AS ENUM (
  'PENDING',
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

CREATE TABLE IF NOT EXISTS "orders_view" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "status" order_status NOT NULL,
  "total_price" float NOT NULL,
  "payment_id" uuid,
  "payment_status" payment_status NOT NULL,
  "payment_image_url" varchar,
  "payment_admin_note" varchar,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "order_items_view" (
  "id" uuid PRIMARY KEY,
  "order_id" uuid,
  "product_id" uuid NOT NULL,
  "product_name" varchar NOT NULL,
  "product_price" varchar NOT NULL,
  "product_quantity" integer NOT NULL,
  "product_image_url" varchar,
  "product_description" text NOT NULL,
  "product_category_id" uuid,
  "product_category_name" varchar NOT NULL,
  "note" text,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "order_addresses_view" (
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

CREATE INDEX ON "orders_view" ("status");
CREATE INDEX ON "order_items_view" ("order_id");
CREATE INDEX ON "order_addresses_view" ("order_id");

ALTER TABLE "order_items_view" ADD FOREIGN KEY ("order_id") REFERENCES "orders_view" ("id");
ALTER TABLE "order_addresses_view" ADD FOREIGN KEY ("order_id") REFERENCES "orders_view" ("id");
