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
    name TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE item (
    id UUID PRIMARY KEY,
    shop_id UUID NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    category_id UUID,
    veg BOOLEAN,
    price NUMERIC(10, 2) NOT NULL,
    image VARCHAR(255),
    scan3d_url TEXT,
    is_available BOOLEAN DEFAULT TRUE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    -- CONSTRAINT fk_category_id FOREIGN KEY (category_id) REFERENCES category(id)
);

CREATE TABLE menu_item (
    id UUID PRIMARY KEY,
    menu_id UUID NOT NULL,
    item_id UUID NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE
    -- CONSTRAINT fk_menu FOREIGN KEY (menu_id) REFERENCES menu(id) ON DELETE CASCADE,
    -- CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES item(id) ON DELETE CASCADE
);

-- Trigger function to update menu_item.active when menu.active changes
CREATE OR REPLACE FUNCTION update_menu_item_active_from_menu()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE menu_item
    SET active = (NEW.active AND i.active)
    FROM item i
    WHERE menu_item.menu_id = NEW.id AND menu_item.item_id = i.id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger on menu table
CREATE TRIGGER trg_update_menu_item_on_menu_update
AFTER UPDATE OF active ON menu
FOR EACH ROW
EXECUTE FUNCTION update_menu_item_active_from_menu();

-- Trigger function to update menu_item.active when item.active changes
CREATE OR REPLACE FUNCTION update_menu_item_active_from_item()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE menu_item
    SET active = (m.active AND NEW.active)
    FROM menu m
    WHERE menu_item.item_id = NEW.id AND menu_item.menu_id = m.id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger on item table
CREATE TRIGGER trg_update_menu_item_on_item_update
AFTER UPDATE OF active ON item
FOR EACH ROW
EXECUTE FUNCTION update_menu_item_active_from_item();

-- Trigger function to update item.isAvailable when category.active changes
CREATE OR REPLACE FUNCTION deactivate_items_on_category_inactive()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE item
        SET is_available = false, category_id = null, updated_at = now()
        WHERE category_id = NEW.id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger on category table
CREATE TRIGGER trg_deactivate_items_on_category_inactive
AFTER UPDATE OF active ON category
FOR EACH ROW
EXECUTE FUNCTION deactivate_items_on_category_inactive();