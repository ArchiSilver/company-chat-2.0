CREATE SCHEMA IF NOT EXISTS messenger;

CREATE TABLE IF NOT EXISTS messenger.chats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(50) NOT NULL DEFAULT 'private', -- 'private' or 'group'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS messenger.chat_members (
    chat_id UUID NOT NULL REFERENCES messenger.chats(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES system.users(id) ON DELETE CASCADE,
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE IF NOT EXISTS messenger.messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id UUID NOT NULL REFERENCES messenger.chats(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES system.users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
