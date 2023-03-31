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


COMMENT ON COLUMN "emails"."id" IS 'email uuid';

COMMENT ON COLUMN "emails"."email" IS 'email address';

COMMENT ON COLUMN "emails"."username" IS 'username of the user';

COMMENT ON COLUMN "emails"."verified" IS 'email verified or not';

COMMENT ON COLUMN "emails"."code" IS 'confirmation code sent to email for email verification';

COMMENT ON COLUMN "emails"."code_expires_at" IS 'email confirmation code expires at';

COMMENT ON COLUMN "emails"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "emails"."updated_at" IS 'last updated at timestamp';

ALTER TABLE "emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
