-- +goose Up
-- +goose StatementBegin
CREATE TABLE "table" (
    id UUID PRIMARY KEY,
    shop_id UUID NOT NULL,
    table_number INT NOT NULL,
    qr_code_url TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "table";
-- +goose StatementEnd
