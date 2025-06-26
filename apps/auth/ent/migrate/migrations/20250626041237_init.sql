-- Create "local_users" table
CREATE TABLE "public"."local_users" (
  "id" uuid NOT NULL,
  "email" character varying NOT NULL,
  "password" character varying NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "is_verified" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "local_users_email_key" to table: "local_users"
CREATE UNIQUE INDEX "local_users_email_key" ON "public"."local_users" ("email");
-- Create "oauth_users" table
CREATE TABLE "public"."oauth_users" (
  "id" uuid NOT NULL,
  "email" character varying NOT NULL,
  "provider" character varying NOT NULL,
  "provider_id" character varying NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "is_verified" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "oauthuser_id_provider" to table: "oauth_users"
CREATE UNIQUE INDEX "oauthuser_id_provider" ON "public"."oauth_users" ("id", "provider");
-- Create index "oauthuser_provider_provider_id" to table: "oauth_users"
CREATE UNIQUE INDEX "oauthuser_provider_provider_id" ON "public"."oauth_users" ("provider", "provider_id");
