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

-- modify "events" table
ALTER TABLE "events"
ADD COLUMN "file_events" uuid NULL,
ADD CONSTRAINT "events_files_events" FOREIGN KEY ("file_events") REFERENCES "files" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
