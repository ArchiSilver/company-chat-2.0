-- Down migration: Drop tasks table

DROP TABLE IF EXISTS tasks.tasks;
DROP SCHEMA IF EXISTS tasks CASCADE;
