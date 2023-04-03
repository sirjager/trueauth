CREATE TABLE "emailrecords" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "verified" BOOL NOT NULL DEFAULT FALSE,
  "token" TEXT NOT NULL DEFAULT '',
  "last_token_sent_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

COMMENT ON COLUMN "emailrecords"."id" IS 'email uuid';

COMMENT ON COLUMN "emailrecords"."email" IS 'email address';

COMMENT ON COLUMN "emailrecords"."verified" IS 'email verified or not';

COMMENT ON COLUMN "emailrecords"."token" IS 'confirmation token';

COMMENT ON COLUMN "emailrecords"."last_token_sent_at" IS 'last time verification requested';

COMMENT ON COLUMN "emailrecords"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "emailrecords"."updated_at" IS 'last updated at timestamp';


CREATE TRIGGER trg_update_updated_at BEFORE UPDATE ON "emailrecords"
FOR EACH ROW WHEN (OLD.id = NEW.id) EXECUTE FUNCTION fn_update_timestamp();
