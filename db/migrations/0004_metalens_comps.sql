-- +goose Up
-- Ligas
CREATE TABLE IF NOT EXISTS app.leagues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    game TEXT NOT NULL CHECK (game IN ('val', 'lol')),
    region TEXT NOT NULL,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE
);

-- Composiciones
-- NOTA: league/map/side ahora son NOT NULL con DEFAULT '' para poder usarlas en UNIQUE.
--       Además creamos un fingerprint del JSON para garantizar unicidad sin meter JSONB directo en UNIQUE.

CREATE TABLE IF NOT EXISTS app.comps (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  game       TEXT NOT NULL CHECK (game IN ('val','lol')),
  region     TEXT NOT NULL,
  league     TEXT NOT NULL DEFAULT '',
  patch      TEXT NOT NULL,
  map        TEXT NOT NULL DEFAULT '',
  side       TEXT NOT NULL DEFAULT '',
  slots      JSONB NOT NULL,
  -- huella del JSON para unique (texto del json -> md5)
  slots_fp   TEXT GENERATED ALWAYS AS (md5(slots::text)) STORED,

  pick_rate  NUMERIC(6,3),
  win_rate   NUMERIC(6,3),
  delta_win  NUMERIC(6,3),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

-- deduplicación lógica
UNIQUE (game, region, league, patch, map, side, slots_fp) );

-- Índices habituales
CREATE INDEX IF NOT EXISTS idx_leagues_game_region ON app.leagues (game, region);

CREATE INDEX IF NOT EXISTS idx_comps_filter ON app.comps (game, region, league, patch);

CREATE INDEX IF NOT EXISTS idx_comps_created ON app.comps (created_at DESC);
-- búsquedas por contenido del JSON
CREATE INDEX IF NOT EXISTS idx_comps_slots_gin ON app.comps USING GIN (slots jsonb_path_ops);

-- +goose Down
DROP INDEX IF EXISTS idx_comps_slots_gin;

DROP INDEX IF EXISTS idx_comps_created;

DROP INDEX IF EXISTS idx_comps_filter;

DROP INDEX IF EXISTS idx_leagues_game_region;

DROP TABLE IF EXISTS app.comps;

DROP TABLE IF EXISTS app.leagues;