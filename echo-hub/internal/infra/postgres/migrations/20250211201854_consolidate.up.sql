-- create "files" table
CREATE TABLE "files" (
  "id" uuid NOT NULL,
  "filename" character varying NOT NULL,
  "filepath" character varying NOT NULL,
  "filesize" bigint NOT NULL,
  "content_type" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

-- create index "files_filepath_key" to table: "files"
CREATE UNIQUE INDEX "files_filepath_key" ON "files" ("filepath");

-- create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL,
  "first_name" character varying NOT NULL DEFAULT '',
  "last_name" character varying NOT NULL DEFAULT '',
  "email" character varying NOT NULL,
  "password" character varying NOT NULL DEFAULT '',
  "access_token" character varying NOT NULL DEFAULT '',
  "refresh_token" character varying NOT NULL DEFAULT '',
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

-- create index "user_access_token" to table: "users"
CREATE INDEX "user_access_token" ON "users" ("access_token");

-- create index "user_email" to table: "users"
CREATE INDEX "user_email" ON "users" ("email");

-- create index "user_email_password" to table: "users"
CREATE INDEX "user_email_password" ON "users" ("email", "password");

-- create index "user_refresh_token" to table: "users"
CREATE INDEX "user_refresh_token" ON "users" ("refresh_token");

-- create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");

-- create "events" table
CREATE TABLE "events" (
  "id" uuid NOT NULL,
  "type" character varying NOT NULL,
  "sub_type" character varying NOT NULL,
  "status" character varying NOT NULL,
  "result" jsonb NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "file_id" uuid NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "events_files_events" FOREIGN KEY ("file_id") REFERENCES "files" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "events_users_events" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
