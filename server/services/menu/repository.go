package menu

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/priyanshubiswas1/barista/menuservice/pkg"
)

type IRepository interface {
	CreateMenu(ctx context.Context, id, shopId, name string, active bool, createdAt time.Time) error
	GetMenus(ctx context.Context, shopId string) ([]pkg.Menu, error)
	SetMenu(ctx context.Context, menuId string) (bool, error)
	DeleteMenu(ctx context.Context, menuId string) (bool, error)
	CreateItem(ctx context.Context, itemId, shopId, name, description, categoryId string, veg bool, price float32, image string, isAvailable, active bool, createdAt time.Time) error
	GetItem(ctx context.Context, itemId string) (pkg.Item, error)
	GetItems(ctx context.Context, menuId string) ([]pkg.Item, error)
	GetItemsAdmin(ctx context.Context, menuId string) ([]pkg.Item, error)
	UpdateItem(ctx context.Context, itemId, shopId, name, description, categoryId string, veg bool, price float32, image string, isAvailable bool, updatedAt time.Time) (bool, error)
	UpdateItemAvailable(ctx context.Context, itemId, shopId string, isAvailable bool, updatedAt time.Time) (bool, error)
	DeleteItem(ctx context.Context, itemId string) (bool, error)
	CreateMenuItem(ctx context.Context, id, menuId, itemId string) (pkg.MenuItem, error)
	GetMenuItems(ctx context.Context, menuId string) ([]pkg.Item, error)
	GetMenuItemsAdmin(ctx context.Context, menuId string) ([]pkg.Item, error)
	DeleteMenuItem(ctx context.Context, menuId, itemId string) (bool, error)
	DeleteMenuItems(ctx context.Context, menuId string) (bool, error)
	GetCategoryName(ctx context.Context, categoryId, shopId string) (string, error)
	CreateCategory(ctx context.Context, id, shopId, name string) error
	UpdateCategory(ctx context.Context, id, shopId, name string) error
	GetCategories(ctx context.Context, shopId string) ([]pkg.Category, error)
	GetCategoriesWithItem(ctx context.Context, shopId string) ([]pkg.CategoryWithItems, error)
	DeleteCategory(ctx context.Context, id, shopId string) (bool, error)
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
	if rowsAffected == 0 {
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
	if rowsAffected == 0 {
		return false, errors.New("failed to delete")
	}
	return true, nil
}

func (r *postgresRepository) CreateItem(ctx context.Context, itemId, shopId, name, description, categoryId string, veg bool, price float32, image string, isAvailable, active bool, createdAt time.Time) error {
	query := `INSERT INTO item
			(id,shop_id,name,description,category_id,veg,price,image,is_available,active,created_at,updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
			`
	res, err := r.db.ExecContext(ctx, query, itemId, shopId, name, description, categoryId, veg, price, image, isAvailable, active, createdAt, createdAt)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("failed to insert")
	}
	return nil
}

func (r *postgresRepository) GetItem(ctx context.Context, itemId string) (pkg.Item, error) {
	query := `SELECT 
				item.id,
				item.shop_id,
				item.name,
				item.description,
				category.name AS category_name,
				item.veg,
				item.price,
				item.image,
				item.is_available
			FROM 
				item
			LEFT JOIN 
				category ON item.category_id = category.id
			WHERE 
				item.id = $1
			`
	res := r.db.QueryRowContext(ctx, query, itemId)
	item := pkg.Item{}
	if err := res.Scan(&item.ItemId, &item.ShopId, &item.Name, &item.Description, &item.Category, &item.Veg, &item.Price, &item.Image, &item.IsAvailable); err != nil {
		if err == sql.ErrNoRows {
			return pkg.Item{}, errors.New("menu item not found")
		}
		return pkg.Item{}, err
	}
	return item, nil
}

func (r *postgresRepository) GetItems(ctx context.Context, shopId string) ([]pkg.Item, error) {
	query := `SELECT 
				item.id,
				item.shop_id,
				item.name,
				item.description,
				category.name AS category_name,
				item.veg,
				item.price,
				item.image,
				item.is_available
			FROM 
				item
			LEFT JOIN 
				category ON item.category_id = category.id
			WHERE 
				item.shop_id = $1
				AND item.active = true
				AND item.is_available = true
			`
	rows, err := r.db.QueryContext(ctx, query, shopId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []pkg.Item{}
	for rows.Next() {
		item := pkg.Item{}
		if err := rows.Scan(&item.ItemId, &item.ShopId, &item.Name, &item.Description, &item.Category, &item.Veg, &item.Price, &item.Image, &item.IsAvailable); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *postgresRepository) GetItemsAdmin(ctx context.Context, shopId string) ([]pkg.Item, error) {
	query := `SELECT 
				item.id,
				item.shop_id,
				item.name,
				item.description,
				category.name AS category_name,
				item.veg,
				item.price,
				item.image,
				item.is_available
			FROM 
				item
			LEFT JOIN 
				category ON item.category_id = category.id
			WHERE 
				item.shop_id = $1
				AND item.active = true
			`
	rows, err := r.db.QueryContext(ctx, query, shopId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []pkg.Item{}
	for rows.Next() {
		item := pkg.Item{}
		if err := rows.Scan(&item.ItemId, &item.ShopId, &item.Name, &item.Description, &item.Category, &item.Veg, &item.Price, &item.Image, &item.IsAvailable); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *postgresRepository) UpdateItem(ctx context.Context, itemId, shopId, name, description, categoryId string, veg bool, price float32, image string, isAvailable bool, updatedAt time.Time) (bool, error) {
	query := `UPDATE item
			SET shop_id=$1,name=$2,description=$3,category_id=$4,veg=$5,price=$6,image=$7,is_available=$8,updated_at=$9
			WHERE id = $10
			`
	res, err := r.db.ExecContext(ctx, query, shopId, name, description, categoryId, veg, price, image, isAvailable, updatedAt, itemId)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, errors.New("failed to update")
	}
	return true, nil
}

func (r *postgresRepository) UpdateItemAvailable(ctx context.Context, itemId, shopId string, isAvailable bool, updatedAt time.Time) (bool, error) {
	var isCatIdExists bool
	getCatQuery := `SELECT (category_id IS NOT NULL) AS is_not_null FROM item WHERE id=$1 AND shop_id=$2`
	resCat := r.db.QueryRowContext(ctx, getCatQuery, itemId, shopId)
	if err := resCat.Scan(&isCatIdExists); err != nil {
		return false, err
	}
	if !isCatIdExists {
		return false, errors.New("missing category Id")
	}

	query := `UPDATE item
			SET is_available=$1,updated_at=$2
			WHERE id = $3 AND shop_id = $4
			`
	res, err := r.db.ExecContext(ctx, query, isAvailable, updatedAt, itemId, shopId)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, errors.New("failed to update")
	}
	return true, nil
}

func (r *postgresRepository) DeleteItem(ctx context.Context, itemId string) (bool, error) {
	query := `UPDATE item
			SET active=false
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
	if rowsAffected == 0 {
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
	query := `SELECT 
				i.id, 
				i.shop_id, 
				i.name, 
				i.description, 
				category.name AS category_name, 
				i.veg,
				i.price, 
				i.image, 
				i.is_available
			FROM 
				item i
			JOIN 
				menu_item mi ON i.id = mi.item_id
			LEFT JOIN 
				category ON i.category_id = category.id
			WHERE 
				mi.menu_id = $1
				AND i.active = true
				AND i.is_available = true
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
			&item.Veg,
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

func (r *postgresRepository) GetMenuItemsAdmin(ctx context.Context, menuId string) ([]pkg.Item, error) {
	query := `SELECT 
				i.id, 
				i.shop_id, 
				i.name, 
				i.description, 
				category.name AS category_name, 
				i.veg,
				i.price, 
				i.image, 
				i.is_available
			FROM 
				item i
			JOIN 
				menu_item mi ON i.id = mi.item_id
			LEFT JOIN 
				category ON i.category_id = category.id
			WHERE 
				mi.menu_id = $1
				AND i.active = true
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
			&item.Veg,
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
	query := `UPDATE menu_item
				SET active = false
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
	if rowsAffected == 0 {
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
	if rowsAffected == 0 {
		return false, errors.New("failed to delete")
	}
	return true, nil
}

func (r *postgresRepository) GetCategoryName(ctx context.Context, categoryId, shopId string) (string, error) {
	catQuery := `SELECT name FROM category 
				WHERE id = $1 AND shop_id = $2
				`
	var categoryName string
	err := r.db.QueryRowContext(ctx, catQuery, categoryId, shopId).Scan(&categoryName)
	if err == sql.ErrNoRows {
		return "", errors.New("category not found")
	}
	if err != nil {
		return "", err
	}
	return categoryName, nil
}

func (r *postgresRepository) CreateCategory(ctx context.Context, id, shopId, name string) error {
	query := `INSERT INTO category
			(id,shop_id,name)
			VALUES ($1,$2,$3)
			`
	rows, err := r.db.ExecContext(ctx, query, id, shopId, name)
	if err != nil {
		return err
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("failed to insert")
	}
	return nil
}

func (r *postgresRepository) UpdateCategory(ctx context.Context, id, shopId, name string) error {
	log.Println("NAME REPO: ", name)
	query := `UPDATE category
			SET name=$3
			WHERE id=$1 AND shop_id=$2
			`
	rows, err := r.db.ExecContext(ctx, query, id, shopId, name)
	if err != nil {
		return err
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("failed to update")
	}
	return nil
}

func (r *postgresRepository) GetCategories(ctx context.Context, shopId string) ([]pkg.Category, error) {
	query := `SELECT id,shop_id,name
			FROM category
			WHERE shop_id = $1
			`
	rows, err := r.db.QueryContext(ctx, query, shopId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := []pkg.Category{}
	for rows.Next() {
		category := pkg.Category{}
		err := rows.Scan(
			&category.Id,
			&category.ShopId,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *postgresRepository) GetCategoriesWithItem(ctx context.Context, shopId string) ([]pkg.CategoryWithItems, error) {
	query := `
			SELECT
			c.id,
			c.shop_id,
			c.name,
			COALESCE(
				json_agg(
					json_build_object(
						'id', i.id,
						'shopId', i.shop_id,
						'name', i.name,
						'description', i.description,
						'categoryId', i.category_id,
						'veg', i.veg,
						'price', i.price,
						'image', i.image,
						'isAvailable', i.is_available
					)
				) FILTER (WHERE i.id IS NOT NULL),
				'[]'
			) AS items
		FROM
			category c
		LEFT JOIN
			item i ON i.category_id = c.id AND i.active = true
		WHERE
			c.shop_id = $1
		GROUP BY
			c.id, c.shop_id, c.name		
	`
	rows, err := r.db.QueryContext(ctx, query, shopId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := []pkg.CategoryWithItems{}
	for rows.Next() {
		category := pkg.CategoryWithItems{}
		var itemsJSON []byte

		if err := rows.Scan(
			&category.Id,
			&category.ShopId,
			&category.Name,
			&itemsJSON,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(itemsJSON, &category.Items); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *postgresRepository) DeleteCategory(ctx context.Context, id, shopId string) (bool, error) {
	query := `DELETE FROM category
			WHERE id=$1 AND shop_id = $2 
			`

	res, err := r.db.ExecContext(ctx, query, id, shopId)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, errors.New("failed to delete")
	}
	return true, nil
}
