-- name: CreateUser :one
-- noinspection SqlResolve
INSERT INTO users (email, password)
VALUES ($1, hash_password($2))
RETURNING *;

-- name: VerifyPassword :one
-- noinspection SqlResolve
SELECT verify_password($1, password)
FROM users
WHERE email = $2;

-- name: CreateOAuthIdentity :one
-- noinspection SqlResolve
INSERT INTO oauth_identities (user_id, provider, id_token, email, raw_data)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByEmail :one
-- noinspection SqlResolve
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserByID :one
-- noinspection SqlResolve
SELECT *
FROM users
WHERE id = $1;

-- name: GetOAuthIdentityByProviderAndToken :one
-- noinspection SqlResolve
SELECT *
FROM oauth_identities
WHERE provider = $1
  AND id_token = $2;