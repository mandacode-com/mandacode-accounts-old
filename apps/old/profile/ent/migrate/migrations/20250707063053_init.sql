-- Create "profiles" table
CREATE TABLE "public"."profiles" (
  "id" uuid NOT NULL,
  "email" character varying NULL,
  "display_name" character varying NULL,
  "bio" character varying NULL,
  "avatar_url" character varying NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
