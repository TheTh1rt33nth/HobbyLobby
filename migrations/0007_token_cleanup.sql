-- +goose Up
-- +goose StatementBegin
-- Speed up the per-request token lookup (WHERE expiry > NOW()) and make the
-- periodic cleanup DELETE cheap.
CREATE INDEX IF NOT EXISTS idx_tokens_expiry ON tokens (expiry);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION cleanup_expired_tokens() RETURNS void
    LANGUAGE sql AS $$
        DELETE FROM tokens WHERE expiry < NOW();
    $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS cleanup_expired_tokens();
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_tokens_expiry;
-- +goose StatementEnd
