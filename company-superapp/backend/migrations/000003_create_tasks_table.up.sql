CREATE SCHEMA IF NOT EXISTS tasks;

CREATE TABLE IF NOT EXISTS tasks.tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'todo' CHECK (status IN ('todo', 'in_progress', 'done')),
    creator_id UUID NOT NULL REFERENCES system.users(id) ON DELETE CASCADE,
    assignee_id UUID REFERENCES system.users(id) ON DELETE SET NULL,
    due_date TIMESTAMPTZ,
    source_message_id BIGINT REFERENCES messenger.messages(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_assignee ON tasks.tasks(assignee_id);
CREATE INDEX idx_tasks_status ON tasks.tasks(status);
