CREATE TABLE "_users" (
  "id" BYTEA PRIMARY KEY NOT NULL,
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "username" VARCHAR(255) UNIQUE NOT NULL,

	"hash_salt" VARCHAR(255) NOT NULL,
	"hash_pass" VARCHAR(255) NOT NULL,

  "firstname" VARCHAR(255) NOT NULL DEFAULT '',
  "lastname" VARCHAR(255) NOT NULL DEFAULT '',

  "verified" BOOL NOT NULL DEFAULT false,
  "blocked" BOOL NOT NULL DEFAULT false,

  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),

  "token_email_verify" TEXT NOT NULL DEFAULT '',
	"token_password_reset" TEXT NOT NULL DEFAULT '',
	"token_email_change" TEXT NOT NULL DEFAULT '',
	"token_user_deletion" TEXT NOT NULL DEFAULT '',

	"last_email_verify" TIMESTAMPTZ NOT NULL DEFAULT '2000-01-01 00:00:00',
	"last_password_reset" TIMESTAMPTZ NOT NULL DEFAULT '2000-01-01 00:00:00',
	"last_email_change" TIMESTAMPTZ NOT NULL DEFAULT '2000-01-01 00:00:00',
	"last_user_deletion" TIMESTAMPTZ NOT NULL DEFAULT '2000-01-01 00:00:00'
);

CREATE INDEX ON "_users" ("firstname");
CREATE INDEX ON "_users" ("lastname");

CREATE TRIGGER trg_update_updated_at BEFORE UPDATE ON "_users"
FOR EACH ROW WHEN (OLD.id = NEW.id) EXECUTE FUNCTION fn_update_timestamp();

COMMENT ON COLUMN "_users"."id" IS 'account id';

COMMENT ON COLUMN "_users"."email" IS 'unique email address';

COMMENT ON COLUMN "_users"."username" IS 'unique username';

COMMENT ON COLUMN "_users"."hash_salt" IS 'salt used for hashing';

COMMENT ON COLUMN "_users"."hash_pass" IS 'hashed password';

COMMENT ON COLUMN "_users"."firstname" IS 'first name';

COMMENT ON COLUMN "_users"."lastname" IS 'last name';

COMMENT ON COLUMN "_users"."verified" IS 'email verified status';

COMMENT ON COLUMN "_users"."blocked" IS 'account blocked status';

COMMENT ON COLUMN "_users"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "_users"."updated_at" IS 'last updated at timestamp';

COMMENT ON COLUMN "_users"."token_email_verify" IS 'email verification token';

COMMENT ON COLUMN "_users"."token_password_reset" IS 'password reset token';

COMMENT ON COLUMN "_users"."token_email_change" IS 'email change token';

COMMENT ON COLUMN "_users"."token_user_deletion" IS 'user deletion token';

COMMENT ON COLUMN "_users"."last_email_verify" IS 'last email verification timestamp';

COMMENT ON COLUMN "_users"."last_password_reset" IS 'last password reset timestamp';

COMMENT ON COLUMN "_users"."last_email_change" IS 'last email change timestamp';

COMMENT ON COLUMN "_users"."last_user_deletion" IS 'last user deletion timestamp';
