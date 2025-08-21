-- +goose Up
-- +goose StatementBegin
CREATE TABLE menu_item (
    id UUID PRIMARY KEY,
    shop_id UUID NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    category TEXT,
    price NUMERIC(10,2) NOT NULL,
    image TEXT,
    is_available BOOLEAN DEFAULT TRUE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS menu_item;
-- +goose StatementEnd
