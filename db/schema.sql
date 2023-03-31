-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-03-31T07:01:35.393Z

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY NOT NULL,
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "username" VARCHAR(255) UNIQUE NOT NULL,
  "password" VARCHAR(255) NOT NULL,
  "firstname" VARCHAR(255) NOT NULL DEFAULT '',
  "lastname" VARCHAR(255) NOT NULL DEFAULT '',
  "verified" BOOL NOT NULL DEFAULT false,
  "blocked" BOOL NOT NULL DEFAULT false,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT 'now()',
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

CREATE INDEX ON "users" ("firstname");

CREATE INDEX ON "users" ("lastname");

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
