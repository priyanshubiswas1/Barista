-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop (
    id UUID PRIMARY KEY,
    owner_id,
    owner_name TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shop;
-- +goose StatementEnd
