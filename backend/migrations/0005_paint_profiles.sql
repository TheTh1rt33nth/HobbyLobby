-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS paint_profiles (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    target_area VARCHAR(255),      -- e.g. 'Armor', 'Skin', 'Base', 'Weapons'
    is_deleted  BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS paint_steps (
    id                  BIGSERIAL PRIMARY KEY,
    paint_profile_id    BIGINT NOT NULL REFERENCES paint_profiles(id) ON DELETE CASCADE,
    step_order          INT NOT NULL,
    paint_name          VARCHAR(255) NOT NULL,   -- e.g. 'Wraithbone', 'Guilliman Flesh'
    brand               VARCHAR(100),            -- e.g. 'Citadel', 'Vallejo', 'Army Painter', 'Scale75'
    paint_type          VARCHAR(50),             -- e.g. 'base', 'contrast', 'shade', 'layer'
    application_method  VARCHAR(100),            -- e.g. 'brush', 'drybrush', 'wash'
    color_hex           CHAR(7),                 -- optional visual reference, e.g. '#F5D7A3'
    notes               TEXT,                    -- e.g. 'Thin 1:1 with medium', 'Focus on recesses only'
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (paint_profile_id, step_order)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS paint_steps;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS paint_profiles;
-- +goose StatementEnd
