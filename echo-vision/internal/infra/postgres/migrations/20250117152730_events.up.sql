-- create "events" table
CREATE TABLE "events" (
  "id" uuid NOT NULL,
  "type" character varying NOT NULL,
  "sub_type" character varying NOT NULL,
  "status" character varying NOT NULL,
  "payload" jsonb NOT NULL,
  "result" jsonb NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "events_users_events" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
