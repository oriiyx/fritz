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