-- +goose Up
CREATE INDEX IF NOT EXISTS idx_costreams_event_live_lang ON app.co_streams (event_id, is_live, lang);

CREATE INDEX IF NOT EXISTS idx_patch_changes_patch ON app.patch_changes (patch_id);

-- +goose Down
DROP INDEX IF EXISTS idx_patch_changes_patch;

DROP INDEX IF EXISTS idx_costreams_event_live_lang;