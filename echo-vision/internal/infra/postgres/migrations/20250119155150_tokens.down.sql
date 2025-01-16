-- reverse: modify "users" table
ALTER TABLE "users" DROP COLUMN "refresh_token", DROP COLUMN "access_token", DROP COLUMN "password";
