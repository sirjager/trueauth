CREATE TABLE "users" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "username" VARCHAR(255) UNIQUE NOT NULL,
  "password" VARCHAR(255) NOT NULL,
  "firstname" VARCHAR(255) NOT NULL DEFAULT '',
  "lastname" VARCHAR(255) NOT NULL DEFAULT '',
  "email_verified" BOOL NOT NULL DEFAULT false,
  "verify_token" TEXT NOT NULL DEFAULT '',
  "last_verify_sent_at" TIMESTAMPTZ NOT NULL,
  "recovery_token" TEXT NOT NULL DEFAULT '',
  "last_recovery_sent_at" TIMESTAMPTZ NOT NULL,
  "emailchange_token" TEXT NOT NULL DEFAULT '',
  "last_emailchange_sent_at" TIMESTAMPTZ NOT NULL,
  "allowed_ips" TEXT[],
  "allowip_token" TEXT NOT NULL DEFAULT '',
  "last_allowip_sent_at" TIMESTAMPTZ NOT NULL,
  "delete_token" TEXT NOT NULL DEFAULT '',
  "last_delete_sent_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("firstname");

CREATE INDEX ON "users" ("lastname");



COMMENT ON COLUMN "users"."id" IS 'user id';

COMMENT ON COLUMN "users"."email" IS 'unique email address';

COMMENT ON COLUMN "users"."username" IS 'unique username';

COMMENT ON COLUMN "users"."password" IS 'hashed password';

COMMENT ON COLUMN "users"."firstname" IS 'first name';

COMMENT ON COLUMN "users"."lastname" IS 'last name';

COMMENT ON COLUMN "users"."email_verified" IS 'email verified status';

COMMENT ON COLUMN "users"."verify_token" IS 'short lived email verification token';

COMMENT ON COLUMN "users"."last_verify_sent_at" IS 'last verification token sent at timestamp';

COMMENT ON COLUMN "users"."recovery_token" IS 'short lived password recovery token';

COMMENT ON COLUMN "users"."last_recovery_sent_at" IS 'last password recovery token sent at timestamp';

COMMENT ON COLUMN "users"."emailchange_token" IS 'short lived email change token';

COMMENT ON COLUMN "users"."last_emailchange_sent_at" IS 'last change email token sent at timestamp';

COMMENT ON COLUMN "users"."allowed_ips" IS 'list of all allowed ip address to access this row';

COMMENT ON COLUMN "users"."allowip_token" IS 'short lived allowip token for allowing new ipaddress';

COMMENT ON COLUMN "users"."last_allowip_sent_at" IS 'last allow ip token sent at timestamp';

COMMENT ON COLUMN "users"."delete_token" IS 'short lived user deletion token';

COMMENT ON COLUMN "users"."last_delete_sent_at" IS 'last deletion token sent at timestamp';

COMMENT ON COLUMN "users"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "users"."updated_at" IS 'last updated at timestamp';


CREATE TRIGGER trg_update_updated_at BEFORE UPDATE ON "users"
FOR EACH ROW WHEN (OLD.id = NEW.id) EXECUTE FUNCTION fn_update_timestamp();
