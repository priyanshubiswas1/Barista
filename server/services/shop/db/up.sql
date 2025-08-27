CREATE TABLE shop (
    id UUID PRIMARY KEY,
    owner_id UUID,
    shop_name TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);