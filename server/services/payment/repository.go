package payment

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/priyanshubiswas1/barista/paymentservice/pkg"
)

type IRepository interface {
	GetPayment(ctx context.Context, sessionId string) (pkg.Payment, error)
	CreatePayment(ctx context.Context, id, sessionId, paymentMethod, paymentStatus, transactionId string, amount float32, createdAt time.Time) error
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

func (r *postgresRepository) GetPayment(ctx context.Context, sessionId string) (pkg.Payment, error) {
	query := `SELECT id,session_id,amount,payment_method,status,transaction_id
				FROM payment WHERE session_id = $1
			`
	row := r.db.QueryRowContext(ctx, query, sessionId)
	payment := pkg.Payment{}
	if err := row.Scan(&payment.Id,
		&payment.SessionId,
		&payment.Amount,
		&payment.PaymentMethod,
		&payment.Status,
		&payment.TransactionId,
	); err != nil {
		return pkg.Payment{}, err
	}
	return payment, nil
}

func (r *postgresRepository) CreatePayment(ctx context.Context, id, sessionId, paymentMethod, status, transactionId string, amount float32, createdAt time.Time) error {
	query := `INSERT INTO payment (id,session_id,amount,payment_method,status,transaction_id,active,created_at,updated_at)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
			`
	res, err := r.db.ExecContext(ctx, query, id, sessionId, amount, paymentMethod, status, transactionId, true, createdAt, createdAt)
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
