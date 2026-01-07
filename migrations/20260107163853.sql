-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "name" character varying(200) NOT NULL,
  "email" character varying(200) NOT NULL,
  "password" character varying(150) NOT NULL,
  "email_verified" boolean NULL DEFAULT false,
  "avatar" character varying(500) NULL DEFAULT NULL::character varying,
  "bio" text NULL,
  "role" "public"."user_role" NULL DEFAULT 'user',
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create "api_keys" table
CREATE TABLE "public"."api_keys" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "name" character varying(100) NULL,
  "key_hash" character varying(255) NOT NULL,
  "is_active" boolean NULL DEFAULT true,
  "last_used_at" timestamp NULL,
  "created_at" timestamp NULL,
  "updated_at" timestamp NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_api_keys_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_api_keys_key_hash" to table: "api_keys"
CREATE UNIQUE INDEX "idx_api_keys_key_hash" ON "public"."api_keys" ("key_hash");
-- Create index "idx_api_keys_user_id" to table: "api_keys"
CREATE INDEX "idx_api_keys_user_id" ON "public"."api_keys" ("user_id");
-- Create "email_verifications" table
CREATE TABLE "public"."email_verifications" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "token" character varying(500) NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "is_verified" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_email_verifications_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_email_verifications_deleted_at" to table: "email_verifications"
CREATE INDEX "idx_email_verifications_deleted_at" ON "public"."email_verifications" ("deleted_at");
-- Create index "idx_email_verifications_token" to table: "email_verifications"
CREATE UNIQUE INDEX "idx_email_verifications_token" ON "public"."email_verifications" ("token");
-- Create index "idx_email_verifications_user_id" to table: "email_verifications"
CREATE INDEX "idx_email_verifications_user_id" ON "public"."email_verifications" ("user_id");
-- Create "password_resets" table
CREATE TABLE "public"."password_resets" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "token" character varying(500) NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "is_used" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_password_resets_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_password_resets_deleted_at" to table: "password_resets"
CREATE INDEX "idx_password_resets_deleted_at" ON "public"."password_resets" ("deleted_at");
-- Create index "idx_password_resets_token" to table: "password_resets"
CREATE UNIQUE INDEX "idx_password_resets_token" ON "public"."password_resets" ("token");
-- Create index "idx_password_resets_user_id" to table: "password_resets"
CREATE INDEX "idx_password_resets_user_id" ON "public"."password_resets" ("user_id");
-- Create "posts" table
CREATE TABLE "public"."posts" (
  "id" bigserial NOT NULL,
  "tweet" character varying(500) NULL,
  "photo" character varying(150) NULL,
  "user_id" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_posts_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_posts_deleted_at" to table: "posts"
CREATE INDEX "idx_posts_deleted_at" ON "public"."posts" ("deleted_at");
-- Create "refresh_tokens" table
CREATE TABLE "public"."refresh_tokens" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "token" character varying(500) NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "is_revoked" boolean NULL DEFAULT false,
  "device_id" character varying(255) NULL,
  "ip_address" character varying(45) NULL,
  "user_agent" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_refresh_tokens_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_refresh_tokens_deleted_at" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_deleted_at" ON "public"."refresh_tokens" ("deleted_at");
-- Create index "idx_refresh_tokens_token" to table: "refresh_tokens"
CREATE UNIQUE INDEX "idx_refresh_tokens_token" ON "public"."refresh_tokens" ("token");
-- Create index "idx_refresh_tokens_user_id" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_user_id" ON "public"."refresh_tokens" ("user_id");
-- Create "user_preferences" table
CREATE TABLE "public"."user_preferences" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "data" jsonb NULL,
  "deleted_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_preferences_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_user_preferences_deleted_at" to table: "user_preferences"
CREATE INDEX "idx_user_preferences_deleted_at" ON "public"."user_preferences" ("deleted_at", "deleted_at");
-- Create index "idx_user_preferences_user_id" to table: "user_preferences"
CREATE UNIQUE INDEX "idx_user_preferences_user_id" ON "public"."user_preferences" ("user_id");
