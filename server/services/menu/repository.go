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
	DeleteMenu(ctx context.Context, menuId string) (bool, error)
	CreateMenuItem(ctx context.Context, itemId, menuId, name, description, category string, price float32, image string, isAvailable, active bool, createdAt time.Time) error
	GetMenuItems(ctx context.Context, menuId string) ([]pkg.MenuItem, error)
	GetMenuItem(ctx context.Context, menuItemId string) (pkg.MenuItem, error)
	UpdateMenuItem(ctx context.Context, itemId, menuId, name, description, category string, price float32, image string, isAvailable, active bool, updatedAt time.Time) (bool, error)
	DeleteMenuItem(ctx context.Context, menuItemId string) (bool, error)
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
	res, err := r.db.QueryContext(ctx, query, shopId)
	if err != nil {
		return nil, err
	}
	menus := []pkg.Menu{}
	for res.Next() {
		menu := pkg.Menu{}
		if err := res.Scan(&menu.MenuId, &menu.ShopId, &menu.Name); err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}
	return menus, nil
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

func (r *postgresRepository) CreateMenuItem(ctx context.Context, itemId, menuId, name, description, category string, price float32, image string,
	isAvailable, active bool, createdAt time.Time) error {

	query := `INSERT INTO menu_item
			(id,menu_id,name,description,category,price,image,is_available,active,created_at,updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
			`
	res, err := r.db.ExecContext(ctx, query, itemId, menuId, name, description, category, price, image, isAvailable, active, createdAt, createdAt)
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

func (r *postgresRepository) GetMenuItems(ctx context.Context, menuId string) ([]pkg.MenuItem, error) {
	query := `SELECT id,menu_id,name,description,category,price,image,is_available
			FROM menu_item
			WHERE menu_id = $1
			`
	res, err := r.db.QueryContext(ctx, query, menuId)
	if err != nil {
		return nil, err
	}
	menuItems := []pkg.MenuItem{}
	for res.Next() {
		menuItem := pkg.MenuItem{}
		if err := res.Scan(&menuItem.MenuItemId, &menuItem.MenuId, &menuItem.Name, &menuItem.Description, &menuItem.Category, &menuItem.Price, &menuItem.Image, &menuItem.IsAvailable); err != nil {
			return nil, err
		}
		menuItems = append(menuItems, menuItem)
	}
	return menuItems, nil
}

func (r *postgresRepository) GetMenuItem(ctx context.Context, menuItemId string) (pkg.MenuItem, error) {
	query := `SELECT id,menu_id,name,description,category,price,image,is_available
			FROM menu_item
			WHERE id = $1
			`
	res := r.db.QueryRowContext(ctx, query, menuItemId)
	menuItem := pkg.MenuItem{}
	if err := res.Scan(&menuItem.MenuItemId, &menuItem.MenuId, &menuItem.Name, &menuItem.Description, &menuItem.Category, &menuItem.Price, &menuItem.Image, &menuItem.IsAvailable); err != nil {
		if err == sql.ErrNoRows {
			return pkg.MenuItem{}, errors.New("menu item not found")
		}
		return pkg.MenuItem{}, err
	}
	return menuItem, nil
}

func (r *postgresRepository) UpdateMenuItem(ctx context.Context, itemId, menuId, name, description, category string, price float32, image string,
	isAvailable, active bool, updatedAt time.Time) (bool, error) {

	query := `UPDATE menu_item
			SET menu_id=$1,name=$2,description=$3,category=$4,price=$5,image=$6,is_available=$7,active=$8,updated_at=$9
			WHERE id = $10
			`
	res, err := r.db.ExecContext(ctx, query, menuId, name, description, category, price, image, isAvailable, active, updatedAt, itemId)
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

func (r *postgresRepository) DeleteMenuItem(ctx context.Context, menuItemId string) (bool, error) {
	query := `DELETE FROM menu_item
			WHERE id = $1
			`
	res, err := r.db.ExecContext(ctx, query, menuItemId)
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
