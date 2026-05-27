-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS project_paint_profiles (
    project_id       BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    paint_profile_id BIGINT NOT NULL REFERENCES paint_profiles(id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, paint_profile_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS project_paint_profiles;
-- +goose StatementEnd
