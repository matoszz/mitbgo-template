-- +goose Up
-- create "todos" table
CREATE TABLE `todos` (`id` text NOT NULL, `name` text NOT NULL, `description` text NULL, PRIMARY KEY (`id`));
-- create index "todo_name" to table: "todos"
CREATE UNIQUE INDEX `todo_name` ON `todos` (`name`);

-- +goose Down
-- reverse: create index "todo_name" to table: "todos"
DROP INDEX `todo_name`;
-- reverse: create "todos" table
DROP TABLE `todos`;
