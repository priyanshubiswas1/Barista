CREATE TABLE session_summary (
    summary_id UUID PRIMARY KEY,
    shop_id UUID NOT NULL,
    session_id UUID NOT NULL,
    session_number INTEGER NOT NULL,
    session_total_amount NUMERIC(10, 2) DEFAULT 0,
    ordered_items JSONB DEFAULT '[]'::jsonb,
    session_closed bool DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
