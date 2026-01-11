-- Down migration: Drop users table and system schema

DROP TABLE IF EXISTS system.users;
DROP SCHEMA IF EXISTS system CASCADE;
