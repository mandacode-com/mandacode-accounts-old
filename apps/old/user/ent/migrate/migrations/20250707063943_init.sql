-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "sync_code" character varying NOT NULL DEFAULT '',
  "email_verification_code" character varying NOT NULL DEFAULT '',
  PRIMARY KEY ("id")
);
