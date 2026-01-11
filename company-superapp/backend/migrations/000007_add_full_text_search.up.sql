-- Add search_vector columns and Full-Text Search support

-- 1. Add search_vector column to users table
ALTER TABLE system.users ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- 2. Add search_vector column to messages table
ALTER TABLE messenger.messages ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- 3. Add search_vector column to tasks table
ALTER TABLE tasks.tasks ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- Create function to update users search_vector
CREATE OR REPLACE FUNCTION system.users_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := 
        setweight(to_tsvector('russian', COALESCE(NEW.full_name, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.full_name, '')), 'A') ||
        setweight(to_tsvector('simple', COALESCE(NEW.email, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create function to update messages search_vector
CREATE OR REPLACE FUNCTION messenger.messages_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := 
        to_tsvector('russian', COALESCE(NEW.content, '')) ||
        to_tsvector('english', COALESCE(NEW.content, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create function to update tasks search_vector
CREATE OR REPLACE FUNCTION tasks.tasks_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := 
        setweight(to_tsvector('russian', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('russian', COALESCE(NEW.description, '')), 'B') ||
        setweight(to_tsvector('english', COALESCE(NEW.description, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
DROP TRIGGER IF EXISTS users_search_vector_trigger ON system.users;
CREATE TRIGGER users_search_vector_trigger
    BEFORE INSERT OR UPDATE ON system.users
    FOR EACH ROW EXECUTE FUNCTION system.users_search_vector_update();

DROP TRIGGER IF EXISTS messages_search_vector_trigger ON messenger.messages;
CREATE TRIGGER messages_search_vector_trigger
    BEFORE INSERT OR UPDATE ON messenger.messages
    FOR EACH ROW EXECUTE FUNCTION messenger.messages_search_vector_update();

DROP TRIGGER IF EXISTS tasks_search_vector_trigger ON tasks.tasks;
CREATE TRIGGER tasks_search_vector_trigger
    BEFORE INSERT OR UPDATE ON tasks.tasks
    FOR EACH ROW EXECUTE FUNCTION tasks.tasks_search_vector_update();

-- Create GIN indexes for fast full-text search
CREATE INDEX IF NOT EXISTS idx_users_search_vector ON system.users USING GIN(search_vector);
CREATE INDEX IF NOT EXISTS idx_messages_search_vector ON messenger.messages USING GIN(search_vector);
CREATE INDEX IF NOT EXISTS idx_tasks_search_vector ON tasks.tasks USING GIN(search_vector);

-- Update existing rows to populate search_vector
UPDATE system.users SET full_name = full_name WHERE search_vector IS NULL;
UPDATE messenger.messages SET content = content WHERE search_vector IS NULL;
UPDATE tasks.tasks SET title = title WHERE search_vector IS NULL;
