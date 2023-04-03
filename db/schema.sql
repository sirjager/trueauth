-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-04-02T11:24:15.505Z

CREATE TABLE "users" (
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
  "user_id" UUID NOT NULL,
  "blocked" BOOL NOT NULL DEFAULT 'false',
  "access_token_expires_at" TIMESTAMPTZ NOT NULL,
  "refresh_token_expires_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE "iprecords" (
  "id" UUID PRIMARY KEY NOT NULL,
  "user_id" UUID UNIQUE NOT NULL,
  "allowed_ips" TEXT[],
  "blocked_ips" TEXT[],
  "token" TEXT NOT NULL DEFAULT '',
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE "emailrecords" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "verified" BOOL NOT NULL DEFAULT FALSE,
  "token" TEXT NOT NULL DEFAULT '',
  "last_token_sent_at" TIMESTAMPTZ NOT NULL,
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

COMMENT ON COLUMN "users"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "users"."updated_at" IS 'last updated at timestamp';

COMMENT ON COLUMN "sessions"."id" IS 'refresh token id';

COMMENT ON COLUMN "sessions"."refresh_token" IS 'refresh token';

COMMENT ON COLUMN "sessions"."access_token_id" IS 'access token id';

COMMENT ON COLUMN "sessions"."access_token" IS 'short life access token';

COMMENT ON COLUMN "sessions"."client_ip" IS 'client ip address';

COMMENT ON COLUMN "sessions"."user_agent" IS 'client user agent';

COMMENT ON COLUMN "sessions"."user_id" IS 'user id to whom this session is assigned to';

COMMENT ON COLUMN "sessions"."blocked" IS 'session is blocked or not';

COMMENT ON COLUMN "sessions"."access_token_expires_at" IS 'expiration time of access token';

COMMENT ON COLUMN "sessions"."refresh_token_expires_at" IS 'expiration time of a refresh token';

COMMENT ON COLUMN "sessions"."created_at" IS 'created at timestamp of this session';

COMMENT ON COLUMN "sessions"."updated_at" IS 'last updated at timestamp of this session';

COMMENT ON COLUMN "iprecords"."id" IS 'record uuid';

COMMENT ON COLUMN "iprecords"."user_id" IS 'user uuid';

COMMENT ON COLUMN "iprecords"."allowed_ips" IS 'list of all allowed ip address for this user';

COMMENT ON COLUMN "iprecords"."blocked_ips" IS 'list of all blocked ip address for this user';

COMMENT ON COLUMN "iprecords"."token" IS 'confirmation token';

COMMENT ON COLUMN "iprecords"."created_at" IS 'created at timestamp of this session';

COMMENT ON COLUMN "iprecords"."updated_at" IS 'last updated at timestamp of this session';

COMMENT ON COLUMN "emailrecords"."id" IS 'email uuid';

COMMENT ON COLUMN "emailrecords"."email" IS 'email address';

COMMENT ON COLUMN "emailrecords"."verified" IS 'email verified or not';

COMMENT ON COLUMN "emailrecords"."token" IS 'confirmation token';

COMMENT ON COLUMN "emailrecords"."last_token_sent_at" IS 'last time verification requested';

COMMENT ON COLUMN "emailrecords"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "emailrecords"."updated_at" IS 'last updated at timestamp';

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
