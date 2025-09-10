-- +goose Up
-- co_streams: ayudar a ORDER BY viewers DESC
CREATE INDEX IF NOT EXISTS idx_costreams_event_live_lang_viewers ON app.co_streams (
    event_id,
    is_live,
    lang,
    viewers DESC
);

-- patches: lista por game y fecha/version
CREATE INDEX IF NOT EXISTS idx_patches_game_rel_ver ON app.patches (
    game,
    released_at DESC,
    version DESC
);

-- patch_changes: filtrar por patch y entity_type
CREATE INDEX IF NOT EXISTS idx_patch_changes_patch_type ON app.patch_changes (patch_id, entity_type);

-- comps: filtros más completos (+ mapa/side)
CREATE INDEX IF NOT EXISTS idx_comps_full_filter ON app.comps (
    game,
    region,
    patch,
    league,
    map,
    side
);

-- comps: para consultas por contenido de slots JSONB
CREATE INDEX IF NOT EXISTS idx_comps_slots_gin ON app.comps USING GIN (slots jsonb_path_ops);

-- +goose Down
DROP INDEX IF EXISTS idx_comps_slots_gin;

DROP INDEX IF EXISTS idx_comps_full_filter;

DROP INDEX IF EXISTS idx_patch_changes_patch_type;

DROP INDEX IF EXISTS idx_patches_game_rel_ver;

DROP INDEX IF EXISTS idx_costreams_event_live_lang_viewers;