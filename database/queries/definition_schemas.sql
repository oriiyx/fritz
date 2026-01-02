-- name: CreateDefinitionSchema :one
-- noinspection SqlResolve
INSERT INTO definition_schemas (id, name, description, schema_json)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDefinitionSchemaByID :one
-- noinspection SqlResolve
SELECT *
FROM definition_schemas
WHERE id = $1;

-- name: GetAllDefinitionSchemas :many
-- noinspection SqlResolve
SELECT *
FROM definition_schemas
ORDER BY name ASC;

-- name: UpdateDefinitionSchema :one
-- noinspection SqlResolve
UPDATE definition_schemas
SET name        = $2,
    description = $3,
    schema_json = $4,
    updated_at  = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteDefinitionSchema :exec
-- noinspection SqlResolve
DELETE
FROM definition_schemas
WHERE id = $1;

-- name: DefinitionSchemaExists :one
-- noinspection SqlResolve
SELECT EXISTS(SELECT 1 FROM definition_schemas WHERE id = $1) as exists;

-- name: GetDefinitionSchemasByUpdatedDate :many
-- Get schemas ordered by most recently updated
-- noinspection SqlResolve
SELECT *
FROM definition_schemas
ORDER BY updated_at DESC
LIMIT $1 OFFSET $2;

-- name: CountDefinitionSchemas :one
-- noinspection SqlResolve
SELECT COUNT(*)
FROM definition_schemas;