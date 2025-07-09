-- Create "auth_accounts" table
CREATE TABLE "public"."auth_accounts" (
  "id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "last_login_at" timestamptz NULL,
  "last_failed_login_at" timestamptz NULL,
  "failed_login_attempts" bigint NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);
-- Create "local_auths" table
CREATE TABLE "public"."local_auths" (
  "id" uuid NOT NULL,
  "email" character varying NOT NULL,
  "password" character varying NOT NULL,
  "is_verified" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "last_login_at" timestamptz NULL,
  "last_failed_login_at" timestamptz NULL,
  "failed_login_attempts" bigint NOT NULL DEFAULT 0,
  "auth_account_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "local_auths_auth_accounts_local_auths" FOREIGN KEY ("auth_account_id") REFERENCES "public"."auth_accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "local_auths_email_key" to table: "local_auths"
CREATE UNIQUE INDEX "local_auths_email_key" ON "public"."local_auths" ("email");
-- Create index "localauth_auth_account_id" to table: "local_auths"
CREATE UNIQUE INDEX "localauth_auth_account_id" ON "public"."local_auths" ("auth_account_id");
-- Create index "localauth_email" to table: "local_auths"
CREATE UNIQUE INDEX "localauth_email" ON "public"."local_auths" ("email");
-- Create "oauth_auths" table
CREATE TABLE "public"."oauth_auths" (
  "id" uuid NOT NULL,
  "provider" character varying NOT NULL,
  "provider_id" character varying NOT NULL,
  "is_verified" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "email" character varying NULL,
  "last_login_at" timestamptz NULL,
  "last_failed_login_at" timestamptz NULL,
  "failed_login_attempts" bigint NOT NULL DEFAULT 0,
  "auth_account_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "oauth_auths_auth_accounts_oauth_auths" FOREIGN KEY ("auth_account_id") REFERENCES "public"."auth_accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "oauthauth_auth_account_id_provider" to table: "oauth_auths"
CREATE UNIQUE INDEX "oauthauth_auth_account_id_provider" ON "public"."oauth_auths" ("auth_account_id", "provider");
-- Create index "oauthauth_provider_provider_id" to table: "oauth_auths"
CREATE UNIQUE INDEX "oauthauth_provider_provider_id" ON "public"."oauth_auths" ("provider", "provider_id");
