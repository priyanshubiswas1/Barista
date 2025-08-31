package menu

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/priyanshubiswas1/barista/menuservice/pkg"
)

type IRepository interface {
	CreateMenu(ctx context.Context, id, shopId, name string, active bool, createdAt time.Time) error
	GetMenus(ctx context.Context, shopId string) ([]pkg.Menu, error)
	SetMenu(ctx context.Context, menuId string) (bool, error)
	DeleteMenu(ctx context.Context, menuId string) (bool, error)
	CreateItem(ctx context.Context, itemId, shopId, name, description, category string, price float32, image string, isAvailable, active bool, createdAt time.Time) error
	GetItem(ctx context.Context, itemId string) (pkg.Item, error)
	GetItems(ctx context.Context, menuId string) ([]pkg.Item, error)
	UpdateItem(ctx context.Context, itemId, shopId, name, description, category string, price float32, image string, isAvailable, active bool, updatedAt time.Time) (bool, error)
	DeleteItem(ctx context.Context, itemId string) (bool, error)
	CreateMenuItem(ctx context.Context, id, menuId, itemId string) (pkg.MenuItem, error)
	GetMenuItems(ctx context.Context, menuId string) ([]pkg.Item, error)
	DeleteMenuItem(ctx context.Context, menuId, itemId string) (bool, error)
	DeleteMenuItems(ctx context.Context, menuId string) (bool, error)
	Close()
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (IRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &postgresRepository{
		db: db,
	}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) CreateMenu(ctx context.Context, id, shopId, name string, active bool, createdAt time.Time) error {
	query := `INSERT INTO menu
			(id,shop_id,name,active,created_at,updated_at)
			VALUES
			($1,$2,$3,$4,$5,$6)
			`
	res, err := r.db.ExecContext(ctx, query, id, shopId, name, active, createdAt, createdAt)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return errors.New("failed to insert")
	}
	return nil
}

func (r *postgresRepository) GetMenus(ctx context.Context, shopId string) ([]pkg.Menu, error) {
	query := `SELECT id,shop_id,name
			FROM menu
			WHERE shop_id = $1
			`
	rows, err := r.db.QueryContext(ctx, query, shopId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	menus := []pkg.Menu{}
	for rows.Next() {
		menu := pkg.Menu{}
		if err := rows.Scan(&menu.MenuId, &menu.ShopId, &menu.Name); err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *postgresRepository) SetMenu(ctx context.Context, menuId string) (bool, error) {
	query := `UPDATE menu
			SET active = CASE
				WHEN id = $1 THEN true
				ELSE false
			END;
			`
	res, err := r.db.ExecContext(ctx, query, menuId)
	if err != nil {
		return false, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *postgresRepository) DeleteMenu(ctx context.Context, menuId string) (bool, error) {
	query := `DELETE FROM menu
			WHERE id = $1
			`
	res, err := r.db.ExecContext(ctx, query, menuId)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected < 1 {
		return false, errors.New("failed to delete")
	}
	return true, nil
}

func (r *postgresRepository) CreateItem(ctx context.Context, itemId, shopId, name, description, category string, price float32, image string, isAvailable, active bool, createdAt time.Time) error {

	query := `INSERT INTO item
			(id,shop_id,name,description,category,price,image,is_available,active,created_at,updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
			`
	res, err := r.db.ExecContext(ctx, query, itemId, shopId, name, description, category, price, image, isAvailable, active, createdAt, createdAt)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return errors.New("failed to insert")
	}
	return nil
}

func (r *postgresRepository) GetItem(ctx context.Context, itemId string) (pkg.Item, error) {
	query := `SELECT id,shop_id,name,description,category,price,image,is_available
			FROM item
			WHERE id = $1
			`
	res := r.db.QueryRowContext(ctx, query, itemId)
	item := pkg.Item{}
	if err := res.Scan(&item.ItemId, &item.ShopId, &item.Name, &item.Description, &item.Category, &item.Price, &item.Image, &item.IsAvailable); err != nil {
		if err == sql.ErrNoRows {
			return pkg.Item{}, errors.New("menu item not found")
		}
		return pkg.Item{}, err
	}
	return item, nil
}

func (r *postgresRepository) GetItems(ctx context.Context, shopId string) ([]pkg.Item, error) {
	query := `SELECT id,shop_id,name,description,category,price,image,is_available
			FROM item
			WHERE shop_id = $1
			`
	rows, err := r.db.QueryContext(ctx, query, shopId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []pkg.Item{}
	for rows.Next() {
		item := pkg.Item{}
		if err := rows.Scan(&item.ItemId, &item.ShopId, &item.Name, &item.Description, &item.Category, &item.Price, &item.Image, &item.IsAvailable); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *postgresRepository) UpdateItem(ctx context.Context, itemId, shopId, name, description, category string, price float32, image string, isAvailable, active bool, updatedAt time.Time) (bool, error) {
	query := `UPDATE item
			SET shop_id=$1,name=$2,description=$3,category=$4,price=$5,image=$6,is_available=$7,active=$8,updated_at=$9
			WHERE id = $10
			`
	res, err := r.db.ExecContext(ctx, query, shopId, name, description, category, price, image, isAvailable, active, updatedAt, itemId)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected < 1 {
		return false, errors.New("failed to update")
	}
	return true, nil
}

func (r *postgresRepository) DeleteItem(ctx context.Context, itemId string) (bool, error) {
	query := `DELETE FROM item
			WHERE id = $1
			`
	res, err := r.db.ExecContext(ctx, query, itemId)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected < 1 {
		return false, errors.New("failed to delete")
	}
	return true, nil
}

func (r *postgresRepository) CreateMenuItem(ctx context.Context, id, menuId, itemId string) (pkg.MenuItem, error) {
	query := `INSERT INTO menu_item
			(id,menu_id,item_id)
			VALUES ($1,$2,$3)
			RETURNING id,menu_id,item_id
			`

	menuItem := pkg.MenuItem{}
	err := r.db.QueryRowContext(ctx, query, id, menuId, itemId).Scan(&menuItem.Id, &menuItem.MenuId, &menuItem.ItemId)
	if err != nil {
		return pkg.MenuItem{}, err
	}
	return menuItem, nil
}

func (r *postgresRepository) GetMenuItems(ctx context.Context, menuId string) ([]pkg.Item, error) {
	query := `SELECT i.id, i.shop_id, i.name, i.description, i.category, i.price, i.image, i.is_available
			FROM item i
			JOIN menu_item mi ON i.id = mi.item_id
			WHERE mi.menu_id = $1
			AND i.active = true;
			`
	rows, err := r.db.QueryContext(ctx, query, menuId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []pkg.Item{}
	for rows.Next() {
		item := pkg.Item{}
		err := rows.Scan(
			&item.ItemId,
			&item.ShopId,
			&item.Name,
			&item.Description,
			&item.Category,
			&item.Price,
			&item.Image,
			&item.IsAvailable,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *postgresRepository) DeleteMenuItem(ctx context.Context, menuId, itemId string) (bool, error) {
	query := `DELETE FROM menu_item
			WHERE item_id=$1 AND menu_id = $2 
			`

	res, err := r.db.ExecContext(ctx, query, itemId, menuId)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected < 1 {
		return false, errors.New("failed to delete")
	}
	return true, nil
}

func (r *postgresRepository) DeleteMenuItems(ctx context.Context, menuId string) (bool, error) {
	query := `DELETE FROM menu_item
			WHERE menu_id = $1
			`

	res, err := r.db.ExecContext(ctx, query, menuId)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected < 1 {
		return false, errors.New("failed to delete")
	}
	return true, nil
}
