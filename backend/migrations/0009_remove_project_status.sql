-- +goose Up
-- +goose StatementBegin
ALTER TABLE projects DROP COLUMN IF EXISTS status;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE projects ADD COLUMN IF NOT EXISTS status hobby_status NOT NULL DEFAULT 'unassembled';
-- +goose StatementEnd
