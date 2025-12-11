-- Entities table - stores basic metadata for all entity instances
CREATE TABLE IF NOT EXISTS entities
(
    id           UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    entity_class TEXT        NOT NULL,                  -- References the entity definition (e.g., 'product', 'customer')
    parent_id    UUID        NULL REFERENCES entities (id) ON DELETE CASCADE,
    o_key        TEXT        NOT NULL,                  -- Object key/name
    o_path       TEXT        NOT NULL,                  -- Full hierarchical path (e.g., '/products/electronics/')
    o_type       TEXT        NOT NULL DEFAULT 'object', -- Type: 'object', 'folder', 'variant'
    published    BOOLEAN     NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by   UUID        NULL REFERENCES users (id) ON DELETE SET NULL,
    updated_by   UUID        NULL REFERENCES users (id) ON DELETE SET NULL,

    -- Ensure unique path + key combination
    CONSTRAINT unique_entity_path_key UNIQUE (o_path, o_key),

    -- Ensure key is URL-safe (no spaces, special chars)
    CONSTRAINT valid_o_key CHECK (o_key ~ '^[a-z0-9_-]+$')
);

-- Index for efficient parent-child queries
CREATE INDEX idx_entities_parent_id ON entities (parent_id);

-- Index for path-based queries (finding all entities under a path)
CREATE INDEX idx_entities_o_path ON entities (o_path);

-- Index for filtering by entity class
CREATE INDEX idx_entities_entity_class ON entities (entity_class);

-- Index for published state queries
CREATE INDEX idx_entities_published ON entities (published);

-- Composite index for common query patterns
CREATE INDEX idx_entities_class_published ON entities (entity_class, published);