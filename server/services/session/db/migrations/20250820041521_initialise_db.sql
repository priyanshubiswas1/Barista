-- +goose Up
-- +goose StatementBegin
CREATE TYPE session_status AS ENUM ('active','closed');
CREATE TABLE session (
    id UUID PRIMARY KEY,
    shop_id UUID NOT NULL,
    table_id UUID NOT NULL,
    status session_status NOT NULL DEFAULT 'active',
    total_amount NUMERIC(10,2) DEFAULT 0,
    is_paid BOOLEAN DEFAULT FALSE,
    payment_method TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS session;
DROP TYPE IF EXISTS session_status;
-- +goose StatementEnd
