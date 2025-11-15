CREATE TABLE menu (
    id UUID PRIMARY KEY,
    shop_id UUID NOT NULL,
    name TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE category (
    id UUID PRIMARY KEY,
    shop_id UUID NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE item (
    id UUID PRIMARY KEY,
    shop_id UUID NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    category_id UUID,
    veg BOOLEAN,
    price NUMERIC(10, 2) NOT NULL,
    image TEXT,
    is_available BOOLEAN DEFAULT TRUE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_category_id FOREIGN KEY (category_id) REFERENCES category(id)
);

CREATE TABLE menu_item (
    id UUID PRIMARY KEY,
    menu_id UUID NOT NULL,
    item_id UUID NOT NULL,
    CONSTRAINT fk_menu FOREIGN KEY (menu_id) REFERENCES menu(id) ON DELETE CASCADE,
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES item(id) ON DELETE CASCADE
);