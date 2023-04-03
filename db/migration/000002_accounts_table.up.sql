CREATE TABLE "accounts" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "username" VARCHAR(255) UNIQUE NOT NULL,
  "password" VARCHAR(255) NOT NULL,
  "firstname" VARCHAR(255) NOT NULL DEFAULT '',
  "lastname" VARCHAR(255) NOT NULL DEFAULT '',
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

COMMENT ON COLUMN "accounts"."created_at" IS 'created at timestamp';

COMMENT ON COLUMN "accounts"."updated_at" IS 'last updated at timestamp';


CREATE TRIGGER trg_update_updated_at BEFORE UPDATE ON "accounts"
FOR EACH ROW WHEN (OLD.id = NEW.id) EXECUTE FUNCTION fn_update_timestamp();
