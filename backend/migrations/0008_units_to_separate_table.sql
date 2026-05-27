-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS units (
    id               BIGSERIAL PRIMARY KEY,
    project_id       BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name             VARCHAR(255) NOT NULL,
    quantity         INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    status           hobby_status NOT NULL DEFAULT 'unassembled',
    notes            TEXT,
    paint_profile_id BIGINT REFERENCES paint_profiles(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS units;
-- +goose StatementEnd
