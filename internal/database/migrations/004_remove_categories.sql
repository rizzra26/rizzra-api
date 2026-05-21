DROP TABLE IF EXISTS stack_items;
DROP TABLE IF EXISTS stack_categories;

CREATE TABLE IF NOT EXISTS stack_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    order_index INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

ALTER TABLE projects DROP COLUMN IF EXISTS category;
