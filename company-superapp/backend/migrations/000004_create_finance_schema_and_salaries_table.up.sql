-- Create finance schema
CREATE SCHEMA IF NOT EXISTS finance;

-- Create salaries table
CREATE TABLE finance.salaries (
    user_id UUID PRIMARY KEY REFERENCES system.users(id) ON DELETE CASCADE,
    amount_encrypted BYTEA NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for faster lookups
CREATE INDEX idx_salaries_user_id ON finance.salaries(user_id);
