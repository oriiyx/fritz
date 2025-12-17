-- name: CreateSession :one
-- noinspection SqlResolve
INSERT INTO sessions (user_identity_id,
                      device_id,
                      is_active,
                      expires_at,
                      created_at,
                      last_activity_at)
VALUES ($1, $2, true, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, user_identity_id, device_id, is_active, created_at, expires_at, last_activity_at;

-- name: GetSessionDetailsByID :one
-- noinspection SqlResolve
SELECT s.id,
       s.user_identity_id,
       s.device_id,
       s.is_active,
       s.created_at,
       s.expires_at,
       s.last_activity_at
FROM sessions s
         JOIN users ui ON s.user_identity_id = ui.id
WHERE s.id = $1
LIMIT 1;

-- name: GetSessionByID :one
-- noinspection SqlResolve
SELECT id,
       user_identity_id,
       device_id,
       is_active,
       created_at,
       expires_at,
       last_activity_at
FROM sessions
WHERE id = $1
LIMIT 1;

-- name: UpdateSessionActivity :one
-- noinspection SqlResolve
UPDATE sessions
SET last_activity_at = CURRENT_TIMESTAMP,
    is_active        = $2
WHERE id = $1
RETURNING id, user_identity_id, device_id, is_active, created_at, expires_at, last_activity_at;

-- name: GetActiveSessionsByUserID :many
-- noinspection SqlResolve
SELECT id, user_identity_id, device_id, is_active, created_at, expires_at, last_activity_at
FROM sessions
WHERE user_identity_id = $1
  AND is_active = true
  AND expires_at > CURRENT_TIMESTAMP;

-- name: DeactivateSession :exec
-- noinspection SqlResolve
UPDATE sessions
SET is_active = false
WHERE id = $1;

-- name: CleanupExpiredSessions :exec
-- noinspection SqlResolve
DELETE
FROM sessions
WHERE expires_at < CURRENT_TIMESTAMP
   OR is_active = false;
