DROP INDEX IF EXISTS idx_sessions_user_identity_id;
DROP INDEX IF EXISTS idx_sessions_device_id;
DROP INDEX IF EXISTS idx_sessions_is_active;
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP TABLE IF EXISTS sessions;