-- Drop triggers
DROP TRIGGER IF EXISTS users_search_vector_trigger ON system.users;
DROP TRIGGER IF EXISTS messages_search_vector_trigger ON messenger.messages;
DROP TRIGGER IF EXISTS tasks_search_vector_trigger ON tasks.tasks;

-- Drop functions
DROP FUNCTION IF EXISTS system.users_search_vector_update();
DROP FUNCTION IF EXISTS messenger.messages_search_vector_update();
DROP FUNCTION IF EXISTS tasks.tasks_search_vector_update();

-- Drop indexes
DROP INDEX IF EXISTS system.idx_users_search_vector;
DROP INDEX IF EXISTS messenger.idx_messages_search_vector;
DROP INDEX IF EXISTS tasks.idx_tasks_search_vector;

-- Drop columns
ALTER TABLE system.users DROP COLUMN IF EXISTS search_vector;
ALTER TABLE messenger.messages DROP COLUMN IF EXISTS search_vector;
ALTER TABLE tasks.tasks DROP COLUMN IF EXISTS search_vector;
