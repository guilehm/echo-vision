-- reverse: modify "events" table
ALTER TABLE "events" DROP CONSTRAINT "events_files_events", DROP COLUMN "file_events";
-- reverse: create index "files_filepath_key" to table: "files"
DROP INDEX "files_filepath_key";
-- reverse: create "files" table
DROP TABLE "files";
