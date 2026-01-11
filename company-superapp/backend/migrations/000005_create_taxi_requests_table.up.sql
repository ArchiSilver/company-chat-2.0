-- Create taxi_requests table in finance schema
CREATE TABLE finance.taxi_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES system.users(id) ON DELETE CASCADE,
    receipt_file_key TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_taxi_requests_user_id ON finance.taxi_requests(user_id);
CREATE INDEX idx_taxi_requests_status ON finance.taxi_requests(status);
