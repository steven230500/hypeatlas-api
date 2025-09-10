-- +goose Up
-- Regla directa: plataforma/handle -> evento
CREATE TABLE IF NOT EXISTS app.event_stream_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    platform TEXT NOT NULL, -- 'twitch' | 'youtube'
    handle TEXT NOT NULL, -- login/handle (minúsculas)
    event_slug TEXT NOT NULL, -- referencia a app.events.slug
    note TEXT DEFAULT '',
    UNIQUE (platform, handle)
);

-- Ventanas activas de un evento (para mapear por calendario)
CREATE TABLE IF NOT EXISTS app.event_windows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    event_slug TEXT NOT NULL, -- app.events.slug
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    region TEXT DEFAULT '', -- opcional (EMEA, NA, ...)
    lang TEXT DEFAULT '' -- opcional (es, en, ...)
);

CREATE INDEX IF NOT EXISTS idx_event_windows_active ON app.event_windows (
    event_slug,
    starts_at,
    ends_at
);

-- +goose Down
DROP TABLE IF EXISTS app.event_windows;

DROP TABLE IF EXISTS app.event_stream_rules;