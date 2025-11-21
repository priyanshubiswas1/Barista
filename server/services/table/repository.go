package table

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/priyanshubiswas1/barista/tableservice/pkg"
)

type IRepository interface {
	CreateTable(ctx context.Context, tableId, shopId, tableNumber, qr, status string, active bool, createdAt time.Time) error
	GetTable(ctx context.Context, tableId string) (pkg.Table, error)
	GetTables(ctx context.Context, shopId string) ([]pkg.Table, error)
	UpdateTable(ctx context.Context, tableId, shopId, tableNumber string, updatedAt time.Time) error
	UpdateTableStatus(ctx context.Context, tableId, shopId, status string, updatedAt time.Time) error
	DeleteTable(ctx context.Context, tableId, shopId string) error
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
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &postgresRepository{
		db: db,
	}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) CreateTable(ctx context.Context, tableId, shopId, tableNumber, qr, status string, active bool, createdAt time.Time) error {
	query := `INSERT INTO "table"
			(id,shop_id,table_number,qr_code_url,table_status,active,created_at,updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)		
	`
	res, err := r.db.ExecContext(ctx, query, tableId, shopId, tableNumber, qr, status, active, createdAt, createdAt)
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

func (r *postgresRepository) GetTable(ctx context.Context, tableId string) (pkg.Table, error) {
	query := `SELECT id,shop_id,table_number,table_status,qr_code_url
			FROM "table"
			WHERE id=$1 AND active = true
			`
	table := pkg.Table{}
	if err := r.db.QueryRowContext(ctx, query, tableId).Scan(&table.TableId, &table.ShopId, &table.TableNumber, &table.Status, &table.QR); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.Table{}, nil
		}
		return pkg.Table{}, err
	}
	return pkg.Table{
		TableId:     table.TableId,
		ShopId:      table.ShopId,
		TableNumber: table.TableNumber,
		Status:      table.Status,
		QR:          table.QR,
	}, nil
}

func (r *postgresRepository) GetTables(ctx context.Context, shopId string) ([]pkg.Table, error) {
	query := `SELECT id,shop_id,table_number,table_status,qr_code_url
			FROM "table"
			WHERE shop_id=$1 AND active = true
			`
	rows, err := r.db.QueryContext(ctx, query, shopId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tables := []pkg.Table{}
	for rows.Next() {
		table := pkg.Table{}
		if err := rows.Scan(&table.TableId, &table.ShopId, &table.TableNumber, &table.Status, &table.QR); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func (r *postgresRepository) UpdateTable(ctx context.Context, tableId, shopId, tableNumber string, updatedAt time.Time) error {
	log.Println(tableId, shopId)
	query := `UPDATE "table"
			SET table_number = $1, updated_at = $2
			WHERE id = $3
			AND shop_id = $4
			`
	res, err := r.db.ExecContext(ctx, query, tableNumber, updatedAt, tableId, shopId)
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

func (r *postgresRepository) UpdateTableStatus(ctx context.Context, tableId, shopId, status string, updatedAt time.Time) error {
	query := `UPDATE "table"
			SET table_status = $1, updated_at = $2
			WHERE id = $3
			AND shop_id = $4
			`
	res, err := r.db.ExecContext(ctx, query, status, updatedAt, tableId, shopId)
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

func (r *postgresRepository) DeleteTable(ctx context.Context, tableId, shopId string) error {
	query := `UPDATE "table"
			SET active = false
			WHERE id = $1
			AND shop_id = $2	
			`
	res, err := r.db.ExecContext(ctx, query, tableId, shopId)
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
