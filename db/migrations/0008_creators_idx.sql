-- +goose Up
CREATE INDEX IF NOT EXISTS idx_creators_platform_verified ON app.creators (platform, verified, handle);

-- +goose Down
DROP INDEX IF EXISTS idx_creators_platform_verified;