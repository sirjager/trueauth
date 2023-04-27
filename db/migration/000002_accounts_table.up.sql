CREATE TABLE "accounts" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "username" VARCHAR(255) UNIQUE NOT NULL,
  "password" VARCHAR(255) NOT NULL,
  "firstname" VARCHAR(255) NOT NULL DEFAULT '',
  "lastname" VARCHAR(255) NOT NULL DEFAULT '',
  "email_verified" BOOL NOT NULL DEFAULT false,
  "confirmation_token" TEXT NOT NULL DEFAULT '',
  "last_confirmation_sent_at" TIMESTAMPTZ NOT NULL,
  "recovery_token" TEXT NOT NULL DEFAULT '',
  "last_recovery_sent_at" TIMESTAMPTZ NOT NULL,
  "email_change_token" TEXT NOT NULL DEFAULT '',
  "last_email_change_sent_at" TIMESTAMPTZ NOT NULL,
  "allowed_ips" TEXT[],
  "allow_ip_token" TEXT NOT NULL DEFAULT '',
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);


CREATE INDEX ON "accounts" ("firstname");

CREATE INDEX ON "accounts" ("lastname");


COMMENT ON COLUMN "accounts"."id" IS 'account id';

COMMENT ON COLUMN "accounts"."email" IS 'unique email';

COMMENT ON COLUMN "accounts"."username" IS 'unique username';

COMMENT ON COLUMN "accounts"."password" IS 'hashed password';

COMMENT ON COLUMN "accounts"."firstname" IS 'first name can be empty';

COMMENT ON COLUMN "accounts"."lastname" IS 'last name can be empty';

COMMENT ON COLUMN "accounts"."email_verified" IS 'email verified status';

COMMENT ON COLUMN "accounts"."confirmation_token" IS 'email confirmation token';

COMMENT ON COLUMN "accounts"."last_confirmation_sent_at" IS 'last email confirmation token sent at';

COMMENT ON COLUMN "accounts"."recovery_token" IS 'password recovery token';

COMMENT ON COLUMN "accounts"."last_recovery_sent_at" IS 'last password recovery token sent at';

COMMENT ON COLUMN "accounts"."email_change_token" IS 'email change token';

COMMENT ON COLUMN "accounts"."last_email_change_sent_at" IS 'last email change token sent at';

COMMENT ON COLUMN "accounts"."allowed_ips" IS 'list of all allowed ip address';

COMMENT ON COLUMN "accounts"."allow_ip_token" IS 'allow logins from ip address';

COMMENT ON COLUMN "accounts"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "accounts"."updated_at" IS 'last updated at timestamp';



CREATE TRIGGER trg_update_updated_at BEFORE UPDATE ON "accounts"
FOR EACH ROW WHEN (OLD.id = NEW.id) EXECUTE FUNCTION fn_update_timestamp();
