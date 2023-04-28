
CREATE TABLE "sessions" (
  "id" UUID PRIMARY KEY NOT NULL,
  "refresh_token" TEXT NOT NULL,
  "access_token_id" UUID UNIQUE NOT NULL,
  "access_token" TEXT NOT NULL,
  "client_ip" TEXT NOT NULL,
  "user_agent" TEXT NOT NULL,
  "user_id" UUID NOT NULL,
  "blocked" BOOL NOT NULL DEFAULT false,
  "access_token_expires_at" TIMESTAMPTZ NOT NULL,
  "refresh_token_expires_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);


CREATE INDEX ON "sessions" ("user_id");

CREATE INDEX ON "sessions" ("access_token_id");


COMMENT ON COLUMN "sessions"."id" IS 'refresh token id';

COMMENT ON COLUMN "sessions"."refresh_token" IS 'refresh token';

COMMENT ON COLUMN "sessions"."access_token_id" IS 'access token id';

COMMENT ON COLUMN "sessions"."access_token" IS 'short lived access token';

COMMENT ON COLUMN "sessions"."client_ip" IS 'client ip address';

COMMENT ON COLUMN "sessions"."user_agent" IS 'client user agent';

COMMENT ON COLUMN "sessions"."user_id" IS 'id of the user assigned to this session';

COMMENT ON COLUMN "sessions"."blocked" IS 'session is blocked or not';

COMMENT ON COLUMN "sessions"."access_token_expires_at" IS 'expiration time of access token';

COMMENT ON COLUMN "sessions"."refresh_token_expires_at" IS 'expiration time of a refresh token';

COMMENT ON COLUMN "sessions"."created_at" IS 'created at timestamp of this session';

COMMENT ON COLUMN "sessions"."updated_at" IS 'last updated at timestamp of this session';

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

CREATE TRIGGER trg_update_updated_at BEFORE UPDATE ON "sessions"
FOR EACH ROW WHEN (OLD.id = NEW.id) EXECUTE FUNCTION fn_update_timestamp();
