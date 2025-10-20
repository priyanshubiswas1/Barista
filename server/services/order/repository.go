package order

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/priyanshubiswas1/barista/orderservice/pkg"
)

type IRepository interface {
	GetOrdersBySessionId(ctx context.Context, sessionId string) ([]pkg.OrderItems, error)
	GetOrderById(ctx context.Context, shopId string) ([]pkg.OrderItems, error)
	CreateOrder(ctx context.Context, id, sessionId, specialNotes string, createdAt time.Time, orderItems []pkg.OrderItemsReq) (pkg.CreateOrderRes, error)
	UpdateOrder(ctx context.Context, orderId, orderStatus string, updatedAt time.Time) error
	CancelOrder(ctx context.Context, orderId string, updatedAT time.Time) error
	GetSessionIdByOrderId(ctx context.Context, orderId string) (string, error)
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
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &postgresRepository{
		db: db,
	}, nil
}

func (r *postgresRepository) Close() error {
	return r.db.Close()
}

func (r *postgresRepository) GetOrdersBySessionId(ctx context.Context, sessionId string) ([]pkg.OrderItems, error) {
	query := `SELECT oi.id, oi.order_id, oi.item_id, oi.quantity, oi.price_at_order_time, o.status
				FROM order_item oi
				JOIN "order" o ON oi.order_id = o.id
				WHERE o.session_id = $1;
			`
	rows, err := r.db.QueryContext(ctx, query, sessionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orderItems := []pkg.OrderItems{}
	for rows.Next() {
		orderItem := pkg.OrderItems{}
		if err := rows.Scan(&orderItem.Id,
			&orderItem.OrderId,
			&orderItem.ItemId,
			&orderItem.Quantity,
			&orderItem.PriceAtOrderTime,
			&orderItem.OrderStatus); err != nil {
			return nil, err
		}
		orderItems = append(orderItems, orderItem)
	}
	return orderItems, nil
}

func (r *postgresRepository) GetOrderById(ctx context.Context, orderId string) ([]pkg.OrderItems, error) {
	query := `SELECT oi.id,oi.order_id,oi.item_id,oi.quantity,oi.price_at_order_time,o.status
				FROM order_item oi 
				JOIN "order" o ON oi.order_id = o.id
				WHERE o.id = $1
			`
	rows, err := r.db.QueryContext(ctx, query, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orderItems := []pkg.OrderItems{}
	for rows.Next() {
		orderItem := pkg.OrderItems{}
		if err := rows.Scan(&orderItem.Id,
			&orderItem.OrderId,
			&orderItem.ItemId,
			&orderItem.Quantity,
			&orderItem.PriceAtOrderTime,
			&orderItem.OrderStatus); err != nil {
			return nil, err
		}
		orderItems = append(orderItems, orderItem)
	}
	return orderItems, nil
}

func (r *postgresRepository) CreateOrder(ctx context.Context, id, sessionId, specialNotes string, createdAt time.Time, orderItems []pkg.OrderItemsReq) (pkg.CreateOrderRes, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return pkg.CreateOrderRes{}, err
	}
	var txErr error
	defer func() {
		if txErr != nil {
			_ = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}
	}()
	query := `
		WITH last_order AS (
			SELECT MAX(order_number) AS max_order_number
			FROM "order"
			WHERE session_id = $2 AND created_at::DATE = CURRENT_DATE
		)
		INSERT INTO "order" (id, session_id, order_number, status, special_notes, active, created_at, updated_at)
		VALUES ($1, $2, COALESCE((SELECT max_order_number + 1 FROM last_order), 1), $3, $4, $5, $6, $7)
	`
	res, err := tx.ExecContext(ctx, query, id, sessionId, "pending", specialNotes, true, createdAt, createdAt)
	if err != nil {
		txErr = err
		return pkg.CreateOrderRes{}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		txErr = err
		return pkg.CreateOrderRes{}, err
	}
	if rowsAffected < 1 {
		txErr = errors.New("failed to create order")
		return pkg.CreateOrderRes{}, txErr
	}
	newOrderItems := []pkg.OrderItems{}
	for _, oi := range orderItems {
		oiid := uuid.New().String()
		query = `INSERT INTO order_item (id,order_id,item_id,quantity,price_at_order_time,active,created_at,updated_at)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		`
		res, err := tx.ExecContext(ctx, query, oiid, id, oi.ItemId, oi.Quantity, oi.PriceAtOrderTime, true, createdAt, createdAt)
		if err != nil {
			txErr = err
			return pkg.CreateOrderRes{}, err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			txErr = err
			return pkg.CreateOrderRes{}, err
		}
		if rowsAffected < 1 {
			txErr = errors.New("failed to create order item")
			return pkg.CreateOrderRes{}, txErr
		}
		newOrderItems = append(newOrderItems, pkg.OrderItems{
			Id:               oiid,
			OrderId:          id,
			ItemId:           oi.ItemId,
			Quantity:         oi.Quantity,
			PriceAtOrderTime: oi.PriceAtOrderTime,
			OrderStatus:      "pending",
		})
	}

	return pkg.CreateOrderRes{
		Order: pkg.Order{
			Id:           id,
			SessionId:    sessionId,
			OrderStatus:  "pending",
			SpecialNotes: specialNotes,
		},
		OrderItems: newOrderItems,
	}, nil
}

func (r *postgresRepository) UpdateOrder(ctx context.Context, orderId, orderStatus string, updatedAT time.Time) error {
	query := `UPDATE "order" SET status = $1,updated_at = $2 WHERE id = $3
			`
	res, err := r.db.ExecContext(ctx, query, orderStatus, updatedAT, orderId)
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

func (r *postgresRepository) CancelOrder(ctx context.Context, orderId string, updatedAt time.Time) error {
	query := `UPDATE "order" SET status = 'cancelled', updated_at = $1 WHERE id = $2
			`
	res, err := r.db.ExecContext(ctx, query, updatedAt, orderId)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return errors.New("failed to cancel order")
	}
	return nil
}

func (r *postgresRepository) GetSessionIdByOrderId(ctx context.Context, orderId string) (string, error) {
	query := `SELECT session_id FROM "order" WHERE id = $1
			`
	var sessionId string
	row := r.db.QueryRowContext(ctx, query, orderId)
	if err := row.Scan(&sessionId); err != nil {
		return "", err
	}
	return sessionId, nil
}
