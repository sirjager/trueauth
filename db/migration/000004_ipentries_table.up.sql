CREATE TABLE "ipentries" (
  "id" uuid PRIMARY KEY NOT NULL,
  "allowed_ips" TEXT[],
  "blocked_ips" TEXT[],
  "code" TEXT NOT NULL DEFAULT '',
  "code_expires_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT 'now()',
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);


COMMENT ON COLUMN "ipentries"."id" IS 'user id';

COMMENT ON COLUMN "ipentries"."allowed_ips" IS 'list of all allowed ip address for this user';

COMMENT ON COLUMN "ipentries"."blocked_ips" IS 'list of all blocked ip address for this user';

COMMENT ON COLUMN "ipentries"."code" IS 'confirmation code sent to email for allowing new ips';

COMMENT ON COLUMN "ipentries"."code_expires_at" IS 'confirmation code expires at';

COMMENT ON COLUMN "ipentries"."created_at" IS 'created at timestamp of this session';

COMMENT ON COLUMN "ipentries"."updated_at" IS 'last updated at timestamp of this session';


ALTER TABLE "ipentries" ADD FOREIGN KEY ("id") REFERENCES "users" ("id");
