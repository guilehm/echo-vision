-- modify "users" table
ALTER TABLE "users"
ADD COLUMN "password" character varying NOT NULL DEFAULT '',
ADD COLUMN "access_token" character varying NOT NULL DEFAULT '',
ADD COLUMN "refresh_token" character varying NOT NULL DEFAULT '';
