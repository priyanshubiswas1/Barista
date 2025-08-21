-- +goose Up
-- +goose StatementBegin
CREATE TYPE payment_status AS ENUM ('pending','completed','failed');
CREATE TABLE payment (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    payment_method TEXT NOT NULL,
    status payment_status NOT NULL DEFAULT 'pending',
    transaction_id TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment;
DROP TYPE IF EXISTS payment_status;
-- +goose StatementEnd
