CREATE TYPE table_status AS ENUM ('free','occupied');
CREATE TABLE "table" (
    id UUID PRIMARY KEY,
    shop_id UUID NOT NULL,
    table_number TEXT NOT NULL,
    qr_code_url TEXT,
    table_status table_status NOT NULL DEFAULT 'free',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);