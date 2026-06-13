CREATE TABLE sessions (
    id            TEXT PRIMARY KEY,
    principal_id  TEXT NOT NULL,
    ip_address    TEXT NOT NULL DEFAULT '',
    user_agent    TEXT NOT NULL DEFAULT '',
    expires_at    TIMESTAMPTZ NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_used_at  TIMESTAMPTZ,
    revoked_at    TIMESTAMPTZ,
    metadata      JSONB DEFAULT '{}'::jsonb
);

CREATE INDEX idx_sessions_principal_id ON sessions(principal_id);
