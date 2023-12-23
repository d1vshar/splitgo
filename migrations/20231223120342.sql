-- Create "_categories" table
CREATE TABLE "_categories" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "created_by" text NULL,
  "updated_by" text NULL,
  "name" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_public__categories_deleted_at" to table: "_categories"
CREATE INDEX "idx_public__categories_deleted_at" ON "_categories" ("deleted_at");
-- Create "_trxs" table
CREATE TABLE "_trxs" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "created_by" text NULL,
  "updated_by" text NULL,
  "description" text NULL,
  "category_id" uuid NULL,
  "amount" numeric NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_public__trxs_category" FOREIGN KEY ("category_id") REFERENCES "_categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_public__trxs_deleted_at" to table: "_trxs"
CREATE INDEX "idx_public__trxs_deleted_at" ON "_trxs" ("deleted_at");
-- Create index "trx_category_idx" to table: "_trxs"
CREATE INDEX "trx_category_idx" ON "_trxs" ("category_id");
-- Create "_trxs_participants" table
CREATE TABLE "_trxs_participants" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "created_by" text NULL,
  "updated_by" text NULL,
  "transaction_id" uuid NOT NULL,
  "user_id" text NOT NULL,
  "amount" numeric NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_public__trxs_transaction_participants" FOREIGN KEY ("transaction_id") REFERENCES "_trxs" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_public__trxs_participants_deleted_at" to table: "_trxs_participants"
CREATE INDEX "idx_public__trxs_participants_deleted_at" ON "_trxs_participants" ("deleted_at");
-- Create index "trx_part_trx_idx" to table: "_trxs_participants"
CREATE INDEX "trx_part_trx_idx" ON "_trxs_participants" ("transaction_id");
-- Create index "trx_part_user_idx" to table: "_trxs_participants"
CREATE INDEX "trx_part_user_idx" ON "_trxs_participants" ("user_id");
