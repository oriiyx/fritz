-- name: CreateEntity :one
-- noinspection SqlResolve
INSERT INTO entities (entity_class,
                      parent_id,
                      o_key,
                      o_path,
                      o_type,
                      published,
                      created_by,
                      updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateEntity :one
-- noinspection SqlResolve
UPDATE entities
SET parent_id  = $1,
    o_key      = $2,
    o_path     = $3,
    published  = $4,
    updated_by = $5,
    updated_at = NOW()
WHERE id = $6
RETURNING *;

-- name: DeleteEntity :exec
-- noinspection SqlResolve
DELETE
FROM entities
WHERE id = $1;

-- name: GetEntityByID :one
-- noinspection SqlResolve
SELECT *
FROM entities
WHERE id = $1;

-- name: GetEntityByPath :one
-- noinspection SqlResolve
SELECT *
FROM entities
WHERE o_path = $1
  AND o_key = $2;