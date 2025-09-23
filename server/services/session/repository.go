package session

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/priyanshubiswas1/barista/sessionservice/pkg"
)

type IRepository interface {
	GetAllSessions(ctx context.Context, shopId string) ([]pkg.Session, error)
	GetActiveSession(ctx context.Context, shopId, tableId string) (pkg.Session, error)
	GetSessionById(ctx context.Context, sessionId string) (pkg.Session, error)
	SendSessionOTP(ctx context.Context) error
	CreateSession(ctx context.Context, id, shopId, tableId, status, paymentMethod, otp string, phoneNumber int32, totalAmount float32, isPaid, active bool, createdAt time.Time) error
	JoinSession(ctx context.Context, shopId, tableId, otp string) (pkg.Session, error)
	UpdateSession(ctx context.Context, id, status, paymentMethod string, totalAmount float32, isPaid, active bool, updatedAt time.Time) error
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

func (r *postgresRepository) GetAllSessions(ctx context.Context, shopId string) ([]pkg.Session, error) {
	query := `SELECT id,shop_id,table_id,otp,status,total_amount,is_paid,payment_method,active,phone_number
				FROM session WHERE shop_id = $1
			`
	rows, err := r.db.QueryContext(ctx, query, shopId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sessions := []pkg.Session{}
	for rows.Next() {
		session := pkg.Session{}
		if err := rows.Scan(&session.Id,
			session.ShopId,
			session.TableId,
			session.Otp,
			session.Status,
			session.Amount,
			session.IsPaid,
			session.PaymentMethod,
			session.Active,
			session.PhoneNumber); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (r *postgresRepository) GetActiveSession(ctx context.Context, shopId, tableId string) (pkg.Session, error) {
	query := `SELECT id,shop_id,table_id,otp,status,total_amount,is_paid,payment_method,active,phone_number
				FROM session WHERE shop_id = $1 AND table_id = $2 AND active = true
			`
	row := r.db.QueryRowContext(ctx, query, shopId, tableId)
	session := pkg.Session{}
	if err := row.Scan(&session.Id,
		session.ShopId,
		session.TableId,
		session.Otp,
		session.Status,
		session.Amount,
		session.IsPaid,
		session.PaymentMethod,
		session.Active,
		session.PhoneNumber); err != nil {
		return pkg.Session{}, err
	}
	return session, nil
}

func (r *postgresRepository) GetSessionById(ctx context.Context, sessionId string) (pkg.Session, error) {
	query := `SELECT id,shop_id,table_id,otp,status,total_amount,is_paid,payment_method,active,phone_number
				FROM session WHERE id = $1
			`
	row := r.db.QueryRowContext(ctx, query, sessionId)
	session := pkg.Session{}
	if err := row.Scan(&session.Id,
		session.ShopId,
		session.TableId,
		session.Otp,
		session.Status,
		session.Amount,
		session.IsPaid,
		session.PaymentMethod,
		session.Active,
		session.PhoneNumber); err != nil {
		return pkg.Session{}, err
	}
	return session, nil
}

func (r *postgresRepository) SendSessionOTP(ctx context.Context) error {
	return nil
}

func (r *postgresRepository) CreateSession(ctx context.Context, id, shopId, tableId, status, paymentMethod, otp string, phoneNumber int32, totalAmount float32, isPaid, active bool, createdAt time.Time) error {
	query := `INSERT INTO session (id,shop_id,table_id,otp,status,total_amount,is_paid,payment_method,active,phone_number,created_at,updated_at)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
			`
	res, err := r.db.ExecContext(ctx, query, id, shopId, tableId, otp, status, totalAmount, isPaid, paymentMethod, active, phoneNumber, createdAt, createdAt)
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

func (r *postgresRepository) JoinSession(ctx context.Context, shopId, tableId, otp string) (pkg.Session, error) {
	query := `SELECT id,shop_id,table_id,otp,status,total_amount,is_paid,payment_method,active,phone_number FROM session WHERE shop_id=$1, table_id=$2`

	row := r.db.QueryRowContext(ctx, query, shopId, tableId)
	session := pkg.Session{}
	if err := row.Scan(&session.Id,
		session.ShopId,
		session.TableId,
		session.Otp,
		session.Status,
		session.Amount,
		session.IsPaid,
		session.PaymentMethod,
		session.Active,
		session.PhoneNumber); err != nil {
		return pkg.Session{}, err
	}
	if session.Otp != otp {
		return pkg.Session{}, errors.New("Invalid OTP")
	}
	return session, nil
}

func (r *postgresRepository) UpdateSession(ctx context.Context, id, status, paymentMethod string, totalAmount float32, isPaid, active bool, updatedAt time.Time) error {
	query := `UPDATE session SET (status=$1,total_amount=$2,is_paid=$3,payment_method=$4,active=$5,updated_at=$6)
				WHERE id=$7
			`
	res, err := r.db.ExecContext(ctx, query, status, totalAmount, isPaid, paymentMethod, active, updatedAt, id)
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
