CREATE TABLE IF NOT EXISTS definition_schemas
(
    id          TEXT PRIMARY KEY,                 -- The definition ID from JSON (e.g., 'product', 'customer')
    name        TEXT        NOT NULL,             -- Human-readable name
    description TEXT,                             -- Optional description
    schema_json JSONB       NOT NULL,             -- The complete definition schema as JSON
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Ensure name is unique as well
    CONSTRAINT unique_definition_name UNIQUE (name)
);

-- Index for faster JSON queries if needed
CREATE INDEX IF NOT EXISTS idx_definition_schemas_json ON definition_schemas USING gin (schema_json);

-- Index for updated_at for efficient sorting/filtering
CREATE INDEX IF NOT EXISTS idx_definition_schemas_updated_at ON definition_schemas (updated_at DESC);