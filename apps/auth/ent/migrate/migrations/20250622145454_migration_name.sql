-- Create "local_users" table
CREATE TABLE "public"."local_users" (
  "id" uuid NOT NULL,
  "email" character varying NOT NULL,
  "password" character varying NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "is_verified" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id")
);
-- Create index "local_users_email_key" to table: "local_users"
CREATE UNIQUE INDEX "local_users_email_key" ON "public"."local_users" ("email");
-- Create "oauth_users" table
CREATE TABLE "public"."oauth_users" (
  "id" uuid NOT NULL,
  "email" character varying NOT NULL,
  "provider" character varying NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "is_verified" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id")
);
