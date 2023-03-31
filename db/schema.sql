-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-03-31T08:00:27.463Z

CREATE TABLE "users" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "username" VARCHAR(255) UNIQUE NOT NULL,
  "password" VARCHAR(255) NOT NULL,
  "firstname" VARCHAR(255) NOT NULL DEFAULT '',
  "lastname" VARCHAR(255) NOT NULL DEFAULT '',
  "verified" BOOL NOT NULL DEFAULT false,
  "blocked" BOOL NOT NULL DEFAULT false,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
  "id" UUID PRIMARY KEY NOT NULL,
  "refresh_token" TEXT NOT NULL,
  "access_token_id" UUID UNIQUE NOT NULL,
  "access_token" TEXT NOT NULL,
  "user_id" UUID NOT NULL,
  "blocked" BOOL NOT NULL DEFAULT 'false',
  "access_token_expires_at" TIMESTAMPTZ NOT NULL,
  "refresh_token_expires_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE "ipentries" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "user_id" UUID UNIQUE NOT NULL,
  "allowed_ips" TEXT[],
  "blocked_ips" TEXT[],
  "code" TEXT NOT NULL DEFAULT '',
  "code_expires_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE "emails" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "email" VARCHAR(255) NOT NULL,
  "username" VARCHAR(255) NOT NULL,
  "verified" BOOL NOT NULL DEFAULT FALSE,
  "code" TEXT NOT NULL DEFAULT '',
  "code_expires_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("firstname");

CREATE INDEX ON "users" ("lastname");

CREATE INDEX ON "sessions" ("user_id");

CREATE INDEX ON "sessions" ("access_token_id");

COMMENT ON COLUMN "users"."id" IS 'user id';

COMMENT ON COLUMN "users"."email" IS 'unique email';

COMMENT ON COLUMN "users"."username" IS 'unique username';

COMMENT ON COLUMN "users"."password" IS 'hashed password';

COMMENT ON COLUMN "users"."firstname" IS 'first name can be empty';

COMMENT ON COLUMN "users"."lastname" IS 'last name can be empty';

COMMENT ON COLUMN "users"."verified" IS 'email verified or not';

COMMENT ON COLUMN "users"."blocked" IS 'user blocked or not';

COMMENT ON COLUMN "users"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "users"."updated_at" IS 'last updated at timestamp';

COMMENT ON COLUMN "sessions"."id" IS 'refresh token id';

COMMENT ON COLUMN "sessions"."refresh_token" IS 'refresh token';

COMMENT ON COLUMN "sessions"."access_token_id" IS 'access token id';

COMMENT ON COLUMN "sessions"."access_token" IS 'short life access token';

COMMENT ON COLUMN "sessions"."user_id" IS 'user id to whom this session is assigned to';

COMMENT ON COLUMN "sessions"."blocked" IS 'session is blocked or not';

COMMENT ON COLUMN "sessions"."access_token_expires_at" IS 'expiration time of access token';

COMMENT ON COLUMN "sessions"."refresh_token_expires_at" IS 'expiration time of a refresh token';

COMMENT ON COLUMN "sessions"."created_at" IS 'created at timestamp of this session';

COMMENT ON COLUMN "sessions"."updated_at" IS 'last updated at timestamp of this session';

COMMENT ON COLUMN "ipentries"."id" IS 'entry uuid';

COMMENT ON COLUMN "ipentries"."user_id" IS 'user id';

COMMENT ON COLUMN "ipentries"."allowed_ips" IS 'list of all allowed ip address for this user';

COMMENT ON COLUMN "ipentries"."blocked_ips" IS 'list of all blocked ip address for this user';

COMMENT ON COLUMN "ipentries"."code" IS 'confirmation code sent to email for allowing new ips';

COMMENT ON COLUMN "ipentries"."code_expires_at" IS 'confirmation code expires at';

COMMENT ON COLUMN "ipentries"."created_at" IS 'created at timestamp of this session';

COMMENT ON COLUMN "ipentries"."updated_at" IS 'last updated at timestamp of this session';

COMMENT ON COLUMN "emails"."id" IS 'email uuid';

COMMENT ON COLUMN "emails"."email" IS 'email address';

COMMENT ON COLUMN "emails"."username" IS 'username of the user';

COMMENT ON COLUMN "emails"."verified" IS 'email verified or not';

COMMENT ON COLUMN "emails"."code" IS 'confirmation code sent to email for email verification';

COMMENT ON COLUMN "emails"."code_expires_at" IS 'email confirmation code expires at';

COMMENT ON COLUMN "emails"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "emails"."updated_at" IS 'last updated at timestamp';

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "ipentries" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
