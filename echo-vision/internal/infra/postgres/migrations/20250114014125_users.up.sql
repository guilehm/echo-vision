-- create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL,
  "first_name" character varying NOT NULL DEFAULT '',
  "last_name" character varying NOT NULL DEFAULT '',
  "email" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

-- create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
