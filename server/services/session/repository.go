package session

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/priyanshubiswas1/barista/sessionservice/pkg"
)

type IRepository interface {
	GetAllSessions(ctx context.Context, shopId string) ([]pkg.Session, error)
	GetActiveSession(ctx context.Context, shopId, tableId string) (pkg.Session, error)
	GetAllActiveSessions(ctx context.Context, shopId string) ([]pkg.Session, error)
	GetSessionById(ctx context.Context, sessionId string) (pkg.Session, error)
	CreateSession(ctx context.Context, id, shopId, tableId, status, paymentMethod, otp, phoneNumber string, totalAmount float32, isPaid, active bool, createdAt time.Time) (string, error)
	JoinSession(ctx context.Context, shopId, tableId, otp string) (pkg.Session, error)
	JoinSessionAdmin(ctx context.Context, shopId, tableId string) (pkg.Session, error)
	UpdateSession(ctx context.Context, id, status, paymentMethod, phoneNumber string, amount float32, isPaid, active bool, updatedAt time.Time) error
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
			&session.ShopId,
			&session.TableId,
			&session.Otp,
			&session.Status,
			&session.Amount,
			&session.IsPaid,
			&session.PaymentMethod,
			&session.Active,
			&session.PhoneNumber); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (r *postgresRepository) GetActiveSession(ctx context.Context, shopId, tableId string) (pkg.Session, error) {
	query := `SELECT id,shop_id,table_id,otp,status,total_amount,is_paid,payment_method,active,phone_number
	FROM session WHERE shop_id = $1 AND table_id = $2 AND status = 'active' AND active = true
	`
	row := r.db.QueryRowContext(ctx, query, shopId, tableId)
	session := pkg.Session{}
	if err := row.Scan(&session.Id,
		&session.ShopId,
		&session.TableId,
		&session.Otp,
		&session.Status,
		&session.Amount,
		&session.IsPaid,
		&session.PaymentMethod,
		&session.Active,
		&session.PhoneNumber); err != nil {
		if err == sql.ErrNoRows {
			return pkg.Session{}, nil
		}
		return pkg.Session{}, err
	}
	return session, nil
}

func (r *postgresRepository) GetAllActiveSessions(ctx context.Context, shopId string) ([]pkg.Session, error) {
	query := `SELECT id,shop_id,table_id,otp,status,total_amount,is_paid,payment_method,active,phone_number
				FROM session WHERE shop_id = $1 AND status = 'active' AND active = true
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
			&session.ShopId,
			&session.TableId,
			&session.Otp,
			&session.Status,
			&session.Amount,
			&session.IsPaid,
			&session.PaymentMethod,
			&session.Active,
			&session.PhoneNumber); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (r *postgresRepository) GetSessionById(ctx context.Context, sessionId string) (pkg.Session, error) {
	query := `SELECT id,shop_id,table_id,otp,status,total_amount,is_paid,payment_method,active,phone_number
	FROM session WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, sessionId)
	session := pkg.Session{}
	if err := row.Scan(&session.Id,
		&session.ShopId,
		&session.TableId,
		&session.Otp,
		&session.Status,
		&session.Amount,
		&session.IsPaid,
		&session.PaymentMethod,
		&session.Active,
		&session.PhoneNumber); err != nil {
		return pkg.Session{}, err
	}
	return session, nil
}

func (r *postgresRepository) CreateSession(ctx context.Context, id, shopId, tableId, status, paymentMethod, otp, phoneNumber string, totalAmount float32, isPaid, active bool, createdAt time.Time) (string, error) {
	var sessionNumber string

	query := `
		WITH last_session AS (
			SELECT MAX(session_number) AS max_session_number
			FROM "session"
			WHERE shop_id = $2 AND created_at::DATE = CURRENT_DATE
		)
		INSERT INTO session (id, session_number, shop_id, table_id, otp, status, total_amount, is_paid, payment_method, active, phone_number, created_at, updated_at)
		VALUES ($1, COALESCE((SELECT max_session_number + 1 FROM last_session), 1), $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING session_number;
	`

	err := r.db.QueryRowContext(ctx, query, id, shopId, tableId, otp, status, totalAmount, isPaid, paymentMethod, active, phoneNumber, createdAt, createdAt).Scan(&sessionNumber)
	if err != nil {
		return "", err
	}

	if sessionNumber == "" {
		return "", errors.New("failed to insert session or retrieve session number")
	}
	return sessionNumber, nil
}

func (r *postgresRepository) JoinSession(ctx context.Context, shopId, tableId, otp string) (pkg.Session, error) {
	query := `SELECT id,shop_id,table_id,otp,status,total_amount,is_paid,payment_method,active,phone_number FROM session WHERE shop_id=$1 AND table_id=$2 AND active=true`

	row := r.db.QueryRowContext(ctx, query, shopId, tableId)
	session := pkg.Session{}
	if err := row.Scan(&session.Id,
		&session.ShopId,
		&session.TableId,
		&session.Otp,
		&session.Status,
		&session.Amount,
		&session.IsPaid,
		&session.PaymentMethod,
		&session.Active,
		&session.PhoneNumber); err != nil {
		return pkg.Session{}, err
	}
	if session.Otp != otp {
		return pkg.Session{}, errors.New("invalid OTP")
	}
	return session, nil
}

func (r *postgresRepository) JoinSessionAdmin(ctx context.Context, shopId, tableId string) (pkg.Session, error) {
	query := `SELECT id,shop_id,table_id,otp,status,total_amount,is_paid,payment_method,active,phone_number FROM session WHERE shop_id=$1 AND table_id=$2 AND active=true`

	row := r.db.QueryRowContext(ctx, query, shopId, tableId)
	session := pkg.Session{}
	if err := row.Scan(&session.Id,
		&session.ShopId,
		&session.TableId,
		&session.Otp,
		&session.Status,
		&session.Amount,
		&session.IsPaid,
		&session.PaymentMethod,
		&session.Active,
		&session.PhoneNumber); err != nil {
		return pkg.Session{}, err
	}
	return session, nil
}

func (r *postgresRepository) UpdateSession(ctx context.Context, id, status, paymentMethod, phoneNumber string, amount float32, isPaid, active bool, updatedAt time.Time) error {
	query := `UPDATE session SET `
	var params []interface{}
	var setClauses []string

	if status != "" {
		setClauses = append(setClauses, "status=$"+strconv.Itoa(len(params)+1))
		params = append(params, status)
	}

	setClauses = append(setClauses, "total_amount=total_amount + $"+strconv.Itoa(len(params)+1))
	params = append(params, amount)

	setClauses = append(setClauses, "is_paid=$"+strconv.Itoa(len(params)+1))
	params = append(params, isPaid)

	if paymentMethod != "" {
		setClauses = append(setClauses, "payment_method=$"+strconv.Itoa(len(params)+1))
		params = append(params, paymentMethod)
	}

	if phoneNumber != "" {
		setClauses = append(setClauses, "phone_number=$"+strconv.Itoa(len(params)+1))
		params = append(params, phoneNumber)
	}

	setClauses = append(setClauses, "active=$"+strconv.Itoa(len(params)+1))
	params = append(params, active)

	setClauses = append(setClauses, "updated_at=$"+strconv.Itoa(len(params)+1))
	params = append(params, updatedAt)

	if len(setClauses) == 0 {
		return errors.New("no values to update")
	}

	query += strings.Join(setClauses, ", ") + " WHERE id=$" + strconv.Itoa(len(params)+1)
	params = append(params, id)

	res, err := r.db.ExecContext(ctx, query, params...)
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

	// query := `UPDATE session SET status=$1,total_amount=total_amount + $2,is_paid=$3,payment_method=$4,active=$5,updated_at=$6
	// 			WHERE id=$7
	// 		`
	// res, err := r.db.ExecContext(ctx, query, status, amount, isPaid, paymentMethod, active, updatedAt, id)
	// if err != nil {
	// 	return err
	// }
	// rowsAffected, err := res.RowsAffected()
	// if err != nil {
	// 	return err
	// }
	// if rowsAffected < 1 {
	// 	return errors.New("failed to update")
	// }
	// return nil
}
