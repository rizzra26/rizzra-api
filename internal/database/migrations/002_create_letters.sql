CREATE TABLE IF NOT EXISTS letters (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug         VARCHAR(255) NOT NULL UNIQUE,
    title        VARCHAR(255) NOT NULL,
    subtitle     TEXT NOT NULL DEFAULT '',
    content      TEXT NOT NULL,
    reading_time INT NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_letters_slug ON letters(slug) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_letters_created_at ON letters(created_at DESC) WHERE deleted_at IS NULL;
