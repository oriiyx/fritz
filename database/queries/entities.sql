-- name: CreateEntity :one
-- noinspection SqlResolve
INSERT INTO entities (entity_class,
                      parent_id,
                      o_key,
                      o_path,
                      o_type,
                      published,
                      created_by,
                      updated_by,
                      has_data)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateEntity :one
-- noinspection SqlResolve
UPDATE entities
SET parent_id  = $1,
    o_key      = $2,
    o_path     = $3,
    published  = $4,
    updated_by = $5,
    has_data   = $6,
    updated_at = NOW()
WHERE id = $7
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

-- name: GetEntityChildren :many
-- Get direct children of a parent entity with pagination
-- noinspection SqlResolve
SELECT e.id,
       e.entity_class,
       e.parent_id,
       e.o_key,
       e.o_path,
       e.o_type,
       e.published,
       e.created_at,
       e.updated_at,
       -- Subquery to check if this entity has children
       EXISTS(SELECT 1
              FROM entities c
              WHERE c.parent_id = e.id
              LIMIT 1)                                            as has_children,
       -- Count of children (useful for UI)
       (SELECT COUNT(*) FROM entities c WHERE c.parent_id = e.id) as children_count
FROM entities e
WHERE e.parent_id = $1
--   AND ($2::text IS NULL OR e.entity_class = $2) -- TODO - potential improvement - Optional class filter
--   AND ($3::text IS NULL OR e.o_type = $3)       -- TODO - potential improvement - Optional type filter
ORDER BY
    -- Folders first, then objects, then variants
    CASE e.o_type
        WHEN 'folder' THEN 1
        WHEN 'object' THEN 2
        WHEN 'variant' THEN 3
        END,
    e.o_key ASC
LIMIT $2 OFFSET $3;

-- name: GetEntityChildrenCount :one
-- Get total count of children (for pagination metadata)
-- noinspection SqlResolve
SELECT COUNT(*)
FROM entities e
WHERE e.parent_id = $1;
--   AND ($2::text IS NULL OR e.entity_class = $2) - todo
--   AND ($3::text IS NULL OR e.o_type = $3); - todo

-- name: GetEntityPath :many
-- Get all ancestors from root to this entity (for breadcrumb)
-- This is the recursive CTE approach
-- noinspection SqlResolve
WITH RECURSIVE entity_path AS (
    -- Start with the target entity
    SELECT id,
           parent_id,
           o_key,
           o_path,
           o_type,
           0 as depth
    FROM entities
    WHERE entities.id = $1

    UNION ALL

    -- Recursively get parents
    SELECT e.id,
           e.parent_id,
           e.o_key,
           e.o_path,
           e.o_type,
           ep.depth + 1
    FROM entities e
             INNER JOIN entity_path ep ON e.id = ep.parent_id)
SELECT id, parent_id, o_key, o_path, o_type, depth
FROM entity_path
ORDER BY depth DESC; -- Root first, target last