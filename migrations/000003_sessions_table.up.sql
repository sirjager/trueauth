CREATE TABLE "_sessions" (
  "id" BYTEA PRIMARY KEY NOT NULL,
  "refresh_token" TEXT NOT NULL,
  "access_token_id" BYTEA UNIQUE NOT NULL,
  "access_token" TEXT NOT NULL,
  "client_ip" TEXT NOT NULL,
  "user_agent" TEXT NOT NULL,
  "user_id" BYTEA NOT NULL,
  "blocked" BOOL NOT NULL DEFAULT false,
  "access_token_expires_at" TIMESTAMPTZ NOT NULL,
  "refresh_token_expires_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

ALTER TABLE "_sessions" ADD FOREIGN KEY ("user_id") REFERENCES "_users" ("id");

CREATE INDEX ON "_sessions" ("user_id");

CREATE INDEX ON "_sessions" ("access_token_id");


COMMENT ON COLUMN "_sessions"."id" IS 'refresh token id';

COMMENT ON COLUMN "_sessions"."refresh_token" IS 'refresh token';

COMMENT ON COLUMN "_sessions"."access_token_id" IS 'access token id';

COMMENT ON COLUMN "_sessions"."access_token" IS 'short lived access token';

COMMENT ON COLUMN "_sessions"."client_ip" IS 'client ip address';

COMMENT ON COLUMN "_sessions"."user_agent" IS 'client user agent';

COMMENT ON COLUMN "_sessions"."user_id" IS 'id of the user assigned to this session';

COMMENT ON COLUMN "_sessions"."blocked" IS 'session is blocked or not';

COMMENT ON COLUMN "_sessions"."access_token_expires_at" IS 'expiration time of access token';

COMMENT ON COLUMN "_sessions"."refresh_token_expires_at" IS 'expiration time of a refresh token';

COMMENT ON COLUMN "_sessions"."created_at" IS 'created at timestamp of this session';

COMMENT ON COLUMN "_sessions"."updated_at" IS 'last updated at timestamp of this session';


CREATE TRIGGER trg_update_updated_at BEFORE UPDATE ON "_sessions"
FOR EACH ROW WHEN (OLD.id = NEW.id) EXECUTE FUNCTION fn_update_timestamp();
