-- Create push_tokens table in system schema
CREATE TABLE system.push_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES system.users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    device_info TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_push_tokens_user_id ON system.push_tokens(user_id);
CREATE UNIQUE INDEX idx_push_tokens_token ON system.push_tokens(token);
