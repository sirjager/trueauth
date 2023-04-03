-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-04-03T13:10:07.579Z

CREATE TABLE "accounts" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "username" VARCHAR(255) UNIQUE NOT NULL,
  "password" VARCHAR(255) NOT NULL,
  "firstname" VARCHAR(255) NOT NULL DEFAULT '',
  "lastname" VARCHAR(255) NOT NULL DEFAULT '',
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
  "id" UUID PRIMARY KEY NOT NULL,
  "refresh_token" TEXT NOT NULL,
  "access_token_id" UUID UNIQUE NOT NULL,
  "access_token" TEXT NOT NULL,
  "client_ip" TEXT NOT NULL,
  "user_agent" TEXT NOT NULL,
  "account_id" UUID NOT NULL,
  "blocked" BOOL NOT NULL DEFAULT 'false',
  "access_token_expires_at" TIMESTAMPTZ NOT NULL,
  "refresh_token_expires_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE "ips" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "account_id" UUID UNIQUE NOT NULL,
  "allowed_ips" TEXT[],
  "blocked_ips" TEXT[],
  "token" TEXT NOT NULL DEFAULT '',
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE "emails" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "verified" BOOL NOT NULL DEFAULT FALSE,
  "token" TEXT NOT NULL DEFAULT '',
  "last_token_sent_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("firstname");

CREATE INDEX ON "accounts" ("lastname");

CREATE INDEX ON "sessions" ("account_id");

CREATE INDEX ON "sessions" ("access_token_id");

COMMENT ON COLUMN "accounts"."id" IS 'account id';

COMMENT ON COLUMN "accounts"."email" IS 'unique email';

COMMENT ON COLUMN "accounts"."username" IS 'unique username';

COMMENT ON COLUMN "accounts"."password" IS 'hashed password';

COMMENT ON COLUMN "accounts"."firstname" IS 'first name can be empty';

COMMENT ON COLUMN "accounts"."lastname" IS 'last name can be empty';

COMMENT ON COLUMN "accounts"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "accounts"."updated_at" IS 'last updated at timestamp';

COMMENT ON COLUMN "sessions"."id" IS 'refresh token id';

COMMENT ON COLUMN "sessions"."refresh_token" IS 'refresh token';

COMMENT ON COLUMN "sessions"."access_token_id" IS 'access token id';

COMMENT ON COLUMN "sessions"."access_token" IS 'short life access token';

COMMENT ON COLUMN "sessions"."client_ip" IS 'client ip address';

COMMENT ON COLUMN "sessions"."user_agent" IS 'client user agent';

COMMENT ON COLUMN "sessions"."account_id" IS 'id of the account assigned to this session';

COMMENT ON COLUMN "sessions"."blocked" IS 'session is blocked or not';

COMMENT ON COLUMN "sessions"."access_token_expires_at" IS 'expiration time of access token';

COMMENT ON COLUMN "sessions"."refresh_token_expires_at" IS 'expiration time of a refresh token';

COMMENT ON COLUMN "sessions"."created_at" IS 'created at timestamp of this session';

COMMENT ON COLUMN "sessions"."updated_at" IS 'last updated at timestamp of this session';

COMMENT ON COLUMN "ips"."id" IS 'record uuid';

COMMENT ON COLUMN "ips"."account_id" IS 'account id';

COMMENT ON COLUMN "ips"."allowed_ips" IS 'list of all allowed ip address';

COMMENT ON COLUMN "ips"."blocked_ips" IS 'list of all blocked ip address';

COMMENT ON COLUMN "ips"."token" IS 'confirmation token';

COMMENT ON COLUMN "ips"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "ips"."updated_at" IS 'last updated at timestamp';

COMMENT ON COLUMN "emails"."id" IS 'email uuid';

COMMENT ON COLUMN "emails"."email" IS 'email address';

COMMENT ON COLUMN "emails"."verified" IS 'email verified or not';

COMMENT ON COLUMN "emails"."token" IS 'confirmation token';

COMMENT ON COLUMN "emails"."last_token_sent_at" IS 'last time verification requested';

COMMENT ON COLUMN "emails"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "emails"."updated_at" IS 'last updated at timestamp';

ALTER TABLE "sessions" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
