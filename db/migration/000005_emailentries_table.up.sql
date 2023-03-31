CREATE TABLE "emailentries" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "user_id" UUID NOT NULL,
  "verified" BOOL NOT NULL DEFAULT FALSE,
  "code" TEXT NOT NULL DEFAULT '',
  "code_expires_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);



COMMENT ON COLUMN "emailentries"."id" IS 'email uuid';

COMMENT ON COLUMN "emailentries"."email" IS 'email address';

COMMENT ON COLUMN "emailentries"."user_id" IS 'user id';

COMMENT ON COLUMN "emailentries"."verified" IS 'email verified or not';

COMMENT ON COLUMN "emailentries"."code" IS 'confirmation code sent to email for email verification';

COMMENT ON COLUMN "emailentries"."code_expires_at" IS 'email confirmation code expires at';

COMMENT ON COLUMN "emailentries"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "emailentries"."updated_at" IS 'last updated at timestamp';

ALTER TABLE "emailentries" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");


CREATE TRIGGER trg_update_updated_at BEFORE UPDATE ON "emailentries"
FOR EACH ROW WHEN (OLD.id = NEW.id) EXECUTE FUNCTION fn_update_timestamp();
