-- +goose Up
ALTER TABLE app.creators
ADD CONSTRAINT uq_creators_platform_handle UNIQUE (platform, handle);

CREATE INDEX IF NOT EXISTS idx_creators_platform_handle ON app.creators (platform, handle);

-- +goose Down
ALTER TABLE app.creators
DROP CONSTRAINT IF EXISTS uq_creators_platform_handle;

DROP INDEX IF EXISTS idx_creators_platform_handle;