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
    has_data     BOOLEAN     NOT NULL DEFAULT false,    -- Tracks if entity has data in its entity_{class} table
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by   UUID        NULL REFERENCES users (id) ON DELETE SET NULL,
    updated_by   UUID        NULL REFERENCES users (id) ON DELETE SET NULL,

    -- Ensure unique path + key combination
    CONSTRAINT unique_entity_path_key UNIQUE (o_path, o_key),

    -- Ensure key is URL-safe (no spaces, special chars)
    CONSTRAINT valid_o_key CHECK (o_key ~ '^[a-zA-Z0-9_-]+$')
);

-- Index for efficient parent-child queries
CREATE INDEX IF NOT EXISTS idx_entities_parent_id ON entities (parent_id);

-- Index for path-based queries (finding all entities under a path)
CREATE INDEX IF NOT EXISTS idx_entities_o_path ON entities (o_path);

-- Index for filtering by entity class
CREATE INDEX IF NOT EXISTS idx_entities_entity_class ON entities (entity_class);

-- Index for published state queries
CREATE INDEX IF NOT EXISTS idx_entities_published ON entities (published);

-- Composite index for common query patterns
CREATE INDEX IF NOT EXISTS idx_entities_class_published ON entities (entity_class, published);

-- Critical for tree navigation queries
CREATE INDEX IF NOT EXISTS idx_entities_tree_nav ON entities (parent_id, entity_class, o_type);

-- Index for querying entities without data (useful for cleanup/validation)
CREATE INDEX IF NOT EXISTS idx_entities_has_data ON entities (has_data);

-- Insert the root entity that all trees start from
INSERT INTO entities (id, entity_class, parent_id, o_key, o_path, o_type, published, has_data)
VALUES ('00000000-0000-0000-0000-000000000001'::uuid, -- Fixed UUID for root
        'system', -- Special class
        NULL, -- No parent
        'system', -- Key
        '/', -- Path is just /
        'folder', -- It's a folder type
        true, -- Always published
        true -- Root always has "data" (it's special)
       )
ON CONFLICT DO NOTHING;