CREATE TABLE "ips" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "account_id" UUID UNIQUE NOT NULL,
  "allowed_ips" TEXT[],
  "blocked_ips" TEXT[],
  "token" TEXT NOT NULL DEFAULT '',
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);



COMMENT ON COLUMN "ips"."id" IS 'record uuid';

COMMENT ON COLUMN "ips"."account_id" IS 'account id';

COMMENT ON COLUMN "ips"."allowed_ips" IS 'list of all allowed ip address';

COMMENT ON COLUMN "ips"."blocked_ips" IS 'list of all blocked ip address';

COMMENT ON COLUMN "ips"."token" IS 'confirmation token';

COMMENT ON COLUMN "ips"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "ips"."updated_at" IS 'last updated at timestamp';


CREATE TRIGGER trg_update_updated_at BEFORE UPDATE ON "ips"
FOR EACH ROW WHEN (OLD.id = NEW.id) EXECUTE FUNCTION fn_update_timestamp();
