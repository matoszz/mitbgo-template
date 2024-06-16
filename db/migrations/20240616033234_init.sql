-- Create "todos" table
CREATE TABLE "todos" ("id" character varying NOT NULL, "name" character varying NOT NULL, "description" character varying NULL, PRIMARY KEY ("id"));
-- Create index "todo_name" to table: "todos"
CREATE UNIQUE INDEX "todo_name" ON "todos" ("name");
