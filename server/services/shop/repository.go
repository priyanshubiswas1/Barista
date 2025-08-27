package shop

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/priyanshubiswas1/shopservice/pkg"
)

type IRepository interface {
	Create(ctx context.Context, ownerId, shopId, shopName string, active bool, createdAt time.Time) error
	GetByShopId(ctx context.Context, shopId string) (*pkg.Shop, error)
	GetByOwnerId(ctx context.Context, ownerId string) ([]*pkg.Shop, error)
	Update(ctx context.Context, ownerId string, shopId, shopName string) error
	Delete(ctx context.Context, ownerId, shopId string) error
	Close() error
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (IRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{
		db: db,
	}, nil
}

func (r *postgresRepository) Close() error {
	return r.db.Close()
}

func (r *postgresRepository) Create(ctx context.Context, ownerId, shopId, shopName string, active bool, createdAt time.Time) error {
	query := `INSERT INTO shop (id,owner_id,shop_name,active,created_at,updated_at)
			VALUES ($1,$2,$3,$4,$5,$6)		
			`
	res, err := r.db.ExecContext(ctx, query, shopId, ownerId, shopName, active, createdAt, createdAt)
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

func (r *postgresRepository) GetByShopId(ctx context.Context, shopId string) (*pkg.Shop, error) {
	query := `SELECT shop_name
			FROM shop
			WHERE id=$1
			`
	var shopName string
	if err := r.db.QueryRowContext(ctx, query, shopId).Scan(&shopName); err != nil {
		return nil, err
	}
	return &pkg.Shop{
		ShopId:   shopId,
		ShopName: shopName,
	}, nil
}

func (r *postgresRepository) GetByOwnerId(ctx context.Context, ownerId string) ([]*pkg.Shop, error) {
	query := `SELECT id,shop_name
			FROM shop
			WHERE owner_id=$1
			`
	rows, err := r.db.QueryContext(ctx, query, ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	shops := []*pkg.Shop{}
	for rows.Next() {
		shop := pkg.Shop{}
		if err := rows.Scan(&shop.ShopId, &shop.ShopName); err != nil {
			return nil, err
		}
		shops = append(shops, &shop)
	}
	return shops, nil
}

func (r *postgresRepository) Update(ctx context.Context, ownerId string, shopId, shopName string) error {
	query := `UPDATE shop
			SET shop_name = $1
			WHERE id = $2
			AND owner_id = $3
			`
	res, err := r.db.ExecContext(ctx, query, shopName, shopId, ownerId)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return errors.New("failed to update")
	}
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, ownerId, shopId string) error {
	query := `DELETE FROM shop
			WHERE id = $1
			AND owner_id = $2	
			`
	res, err := r.db.ExecContext(ctx, query, shopId, ownerId)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return errors.New("failed to delete")
	}
	return nil
}
