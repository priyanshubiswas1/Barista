CREATE TABLE menu (
    id UUID PRIMARY KEY,
    shop_id UUID NOT NULL,
    name TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE menu_item (
    id UUID PRIMARY KEY,
    menu_id UUID NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    category TEXT NOT NULL CHECK (category IN ('Veg', 'NonVeg', 'Dessert')),
    price NUMERIC(10,2) NOT NULL,
    image TEXT,
    is_available BOOLEAN DEFAULT TRUE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);