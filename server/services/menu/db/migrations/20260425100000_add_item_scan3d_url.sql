-- +goose Up
-- +goose StatementBegin
ALTER TABLE item
ADD COLUMN IF NOT EXISTS scan3d_url TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE item
DROP COLUMN IF EXISTS scan3d_url;
-- +goose StatementEnd

