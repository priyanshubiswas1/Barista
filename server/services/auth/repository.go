package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/priyanshubiswas1/authservice/pkg"
)

type Repository interface {
	Close()
	Register(ctx context.Context, id string, user pkg.User, active bool, createdAt time.Time) error
	Login(ctx context.Context) error
	Getuser(ctx context.Context) error
	UpdateUser(ctx context.Context) error
	DeleteUser(ctx context.Context) error
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) Register(ctx context.Context, id string, user pkg.User, active bool, createdAt time.Time) error {
	query := `
			INSERT INTO users
			(id,name,email,password,phone,address,active,created_at,updated_at) 
			VALUES 
			($1,$2,$3,$4,$5,$6,$7,$8,$9)
			`
	res, err := r.db.ExecContext(ctx, query,
		id,
		user.Name,
		user.Email,
		user.Password,
		user.Phone,
		user.Address,
		active,
		createdAt,
		createdAt,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return errors.New("insert failed: rows affected 0")
	}
	return nil
}

func (r *postgresRepository) Login(ctx context.Context) error {
	return nil
}

func (r *postgresRepository) Getuser(ctx context.Context) error {
	return nil
}

func (r *postgresRepository) UpdateUser(ctx context.Context) error {
	return nil
}

func (r *postgresRepository) DeleteUser(ctx context.Context) error {
	return nil
}
