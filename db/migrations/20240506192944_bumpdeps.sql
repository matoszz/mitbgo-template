-- Create "todos" table
CREATE TABLE `todos` (`id` text NOT NULL, `name` text NOT NULL, `description` text NULL, PRIMARY KEY (`id`));
-- Create index "todo_name" to table: "todos"
CREATE UNIQUE INDEX `todo_name` ON `todos` (`name`);
