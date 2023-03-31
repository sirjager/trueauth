CREATE TABLE "iprecords" (
  "id" UUID PRIMARY KEY NOT NULL,
  "allowed_ips" TEXT[],
  "blocked_ips" TEXT[],
  "code" TEXT NOT NULL DEFAULT '',
  "code_expires_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);



COMMENT ON COLUMN "iprecords"."id" IS 'user uuid';

COMMENT ON COLUMN "iprecords"."allowed_ips" IS 'list of all allowed ip address for this user';

COMMENT ON COLUMN "iprecords"."blocked_ips" IS 'list of all blocked ip address for this user';

COMMENT ON COLUMN "iprecords"."code" IS 'confirmation code sent to email for allowing new ips';

COMMENT ON COLUMN "iprecords"."code_expires_at" IS 'confirmation code expires at';

COMMENT ON COLUMN "iprecords"."created_at" IS 'created at timestamp of this session';

COMMENT ON COLUMN "iprecords"."updated_at" IS 'last updated at timestamp of this session';


CREATE TRIGGER trg_update_updated_at BEFORE UPDATE ON "iprecords"
FOR EACH ROW WHEN (OLD.id = NEW.id) EXECUTE FUNCTION fn_update_timestamp();
