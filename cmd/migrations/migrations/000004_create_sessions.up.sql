CREATE TABLE IF NOT EXISTS sessions
(
    id               UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    user_identity_id UUID        NOT NULL REFERENCES users (id),
    device_id        TEXT        NOT NULL,
    is_active        BOOLEAN     NOT NULL DEFAULT true,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_activity_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user_identity_id ON sessions (user_identity_id);
CREATE INDEX idx_sessions_device_id ON sessions (device_id);
CREATE INDEX idx_sessions_is_active ON sessions (is_active);
CREATE INDEX idx_sessions_expires_at ON sessions (expires_at);