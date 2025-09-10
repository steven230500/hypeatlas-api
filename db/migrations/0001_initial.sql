-- +goose Up
CREATE SCHEMA IF NOT EXISTS app;

-- Extensiones útiles (UUID)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =========================
-- RELAY (HypeMap)
-- =========================
CREATE TABLE IF NOT EXISTS app.events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    game TEXT NOT NULL CHECK (game IN ('val', 'lol')),
    league TEXT,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS app.creators (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    platform TEXT NOT NULL CHECK (
        platform IN ('twitch', 'youtube')
    ),
    handle TEXT NOT NULL,
    url TEXT NOT NULL,
    lang TEXT NOT NULL, -- ISO-639-1 (es, en, fr…)
    country TEXT, -- ISO-3166-1 alpha-2 (ES, MX…)
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS app.co_streams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    event_id UUID NOT NULL REFERENCES app.events (id) ON DELETE CASCADE,
    creator_id UUID NOT NULL REFERENCES app.creators (id) ON DELETE CASCADE,
    platform TEXT NOT NULL CHECK (
        platform IN ('twitch', 'youtube')
    ),
    url TEXT NOT NULL,
    lang TEXT NOT NULL,
    country TEXT,
    viewers INTEGER NOT NULL DEFAULT 0,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_live BOOLEAN NOT NULL DEFAULT FALSE,
    last_seen_at TIMESTAMPTZ,
    UNIQUE (event_id, creator_id)
);

-- =========================
-- SIGNAL (MetaLens)
-- =========================
CREATE TABLE IF NOT EXISTS app.patches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    game TEXT NOT NULL CHECK (game IN ('val', 'lol')),
    version TEXT NOT NULL,
    released_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS app.patch_changes (
    id BIGSERIAL PRIMARY KEY,
    patch_id UUID NOT NULL REFERENCES app.patches (id) ON DELETE CASCADE,
    entity_type TEXT NOT NULL, -- champion|agent|item|weapon|map
    entity_id TEXT NOT NULL,
    field TEXT NOT NULL,
    old TEXT,
    new TEXT,
    impact_score NUMERIC(5, 2) NOT NULL DEFAULT 0.00
);

-- +goose Down
DROP TABLE IF EXISTS app.patch_changes;

DROP TABLE IF EXISTS app.patches;

DROP TABLE IF EXISTS app.co_streams;

DROP TABLE IF EXISTS app.creators;

DROP TABLE IF EXISTS app.events;

DROP SCHEMA IF EXISTS app CASCADE;