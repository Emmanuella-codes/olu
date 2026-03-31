CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS candidates (
    id UUID PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    party TEXT NOT NULL,
    bio TEXT NOT NULL DEFAULT '',
    achievements TEXT NOT NULL DEFAULT '',
    photo_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_candidates_is_active_name
    ON candidates (is_active, name);

CREATE TABLE IF NOT EXISTS admins (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_admins_email
    ON admins (email);

CREATE TABLE IF NOT EXISTS votes (
    id UUID PRIMARY KEY,
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE RESTRICT,
    voter_hash TEXT NOT NULL UNIQUE,
    channel TEXT NOT NULL CHECK (channel IN ('sms', 'web')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_votes_candidate_id
    ON votes (candidate_id);

CREATE INDEX IF NOT EXISTS idx_votes_created_at
    ON votes (created_at);

CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    voter_hash TEXT NOT NULL,
    candidate_code TEXT NOT NULL,
    channel TEXT NOT NULL,
    status TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_voter_hash_created_at
    ON audit_logs (voter_hash, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_logs_status_created_at
    ON audit_logs (status, created_at DESC);

CREATE TABLE IF NOT EXISTS sms_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    msisdn TEXT NOT NULL,
    raw_message TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'processed', 'rejected')),
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ,
    rejected_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_sms_queue_status_created_at
    ON sms_queue (status, created_at);
