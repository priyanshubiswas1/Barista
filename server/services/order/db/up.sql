CREATE TYPE order_status AS ENUM ('pending','in_progress','ready','served','cancelled');
CREATE TABLE "order" (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL,
    status order_status NOT NULL DEFAULT 'pending',
    special_notes TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE order_item (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    menu_item_id UUID NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    price_at_order_time NUMERIC(10,2) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);