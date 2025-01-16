-- create index "user_access_token" to table: "users"
CREATE INDEX "user_access_token" ON "users" ("access_token");
-- create index "user_email" to table: "users"
CREATE INDEX "user_email" ON "users" ("email");
-- create index "user_email_password" to table: "users"
CREATE INDEX "user_email_password" ON "users" ("email", "password");
-- create index "user_refresh_token" to table: "users"
CREATE INDEX "user_refresh_token" ON "users" ("refresh_token");
