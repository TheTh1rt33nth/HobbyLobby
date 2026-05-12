-- +goose Up
-- +goose StatementBegin
CREATE TYPE hobby_status AS ENUM (
    'unassembled',
    'assembled',
    'primed',
    'base_coated',
    'painted',
    'based',
    'complete'
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE projects
    ADD COLUMN IF NOT EXISTS game_system VARCHAR(100),
    ADD COLUMN IF NOT EXISTS faction     VARCHAR(100),
    ADD COLUMN IF NOT EXISTS status      hobby_status NOT NULL DEFAULT 'unassembled';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE projects
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS game_system,
    DROP COLUMN IF EXISTS faction;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TYPE IF EXISTS hobby_status;
-- +goose StatementEnd
