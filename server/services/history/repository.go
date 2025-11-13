package history

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/priyanshubiswas1/barista/historyservice/pkg"
	history "github.com/priyanshubiswas1/baristacontracts/gen/history"
)

type IRepository interface {
	GetSessionSummary(ctx context.Context, offset, limit int32, shopId string, timeRange history.TimeRange) ([]pkg.Summary, int32, int32, error)
	CreateSession(ctx context.Context, summaryId, shopId, sessionId string, sessionNumber int32, totalAmount float32, orderedItems []byte, sessionClosed bool, createdAt time.Time) error
	CloseSession(ctx context.Context, sessionId string, sessionClosed bool, updatedAt time.Time) error
	UpdateSession(ctx context.Context, sessionId string, totalAmount float32, updatedAt time.Time) (string, error)
	AddOrder(ctx context.Context, sessionId string, newItems []pkg.OrderedItem) (pkg.Summary, error)
	RemoveOrderItems(ctx context.Context, sessionId string, orderItemIds []string) error
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
	return &postgresRepository{
		db: db,
	}, nil
}

func (s *postgresRepository) GetSessionSummary(ctx context.Context, offset, limit int32, shopId string, timeRange history.TimeRange) ([]pkg.Summary, int32, int32, error) {
	now := time.Now()
	var startDate time.Time

	switch timeRange {
	case history.TimeRange_TODAY:
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case history.TimeRange_LAST_WEEK:
		startDate = now.AddDate(0, 0, -7)
	case history.TimeRange_LAST_MONTH:
		startDate = now.AddDate(0, -1, 0)
	case history.TimeRange_LAST_3_MONTHS:
		startDate = now.AddDate(0, -3, 0)
	default:
		startDate = time.Time{}
	}

	baseQuery := `FROM session_summary
				  WHERE shop_id = $1 AND session_closed = true`
	baseArgs := []interface{}{shopId}
	argIndex := 2

	if !startDate.IsZero() {
		baseQuery += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		baseArgs = append(baseArgs, startDate)
		argIndex++
	}

	countQuery := `SELECT COUNT(*), COALESCE(SUM(session_total_amount), 0) ` + baseQuery
	var totalSessions int32
	var totalAmount float64
	err := s.db.QueryRowContext(ctx, countQuery, baseArgs...).Scan(&totalSessions, &totalAmount)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get totals: %w", err)
	}

	query := `SELECT summary_id, shop_id, session_id, session_number, session_total_amount, ordered_items, created_at, updated_at ` +
		baseQuery + ` ORDER BY created_at DESC`
	
	queryArgs := make([]interface{}, len(baseArgs))
	copy(queryArgs, baseArgs)
	currentArgIndex := argIndex
	
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", currentArgIndex)
		queryArgs = append(queryArgs, limit)
		currentArgIndex++
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", currentArgIndex)
		queryArgs = append(queryArgs, offset)
	}

	rows, err := s.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var summaries []pkg.Summary
	for rows.Next() {
		var summary pkg.Summary
		var orderedItemsJSON []byte

		err := rows.Scan(
			&summary.Id,
			&summary.ShopId,
			&summary.SessionId,
			&summary.SessionNumber,
			&summary.TotalAmount,
			&orderedItemsJSON,
			&summary.CreatedAt,
			&summary.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		if len(orderedItemsJSON) > 0 {
			var orderedItems []pkg.OrderedItem
			err = json.Unmarshal(orderedItemsJSON, &orderedItems)
			if err != nil {
				log.Printf("failed to unmarshal ordered_items for summary %s: %v", summary.Id, err)
				orderedItems = []pkg.OrderedItem{}
			}
			summary.OrderedItems = orderedItems
		}
		summaries = append(summaries, summary)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	return summaries, int32(totalAmount), totalSessions, nil
}

func (s *postgresRepository) CreateSession(ctx context.Context, summaryId, shopId, sessionId string, sessionNumber int32, totalAmount float32, orderedItems []byte, sessionClosed bool, createdAt time.Time) error {
	// Convert nil to empty JSON array for PostgreSQL JSONB
	if orderedItems == nil {
		orderedItems = []byte("[]")
	}
	query := `INSERT INTO session_summary (summary_id, shop_id, session_id, session_number, session_total_amount, ordered_items,session_closed, created_at, updated_at)
				VALUES ($1,$2,$3,$4,$5,$6::jsonb,$7,$8,$9);
			`
	res, err := s.db.ExecContext(ctx, query, summaryId, shopId, sessionId, sessionNumber, totalAmount, orderedItems, sessionClosed, createdAt, createdAt)
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

func (s *postgresRepository) UpdateSession(ctx context.Context, sessionId string, totalAmount float32, updatedAt time.Time) (string, error) {
	query := `
		UPDATE session_summary
		SET session_total_amount = $1, updated_at = $2
		WHERE session_id = $3
		RETURNING summary_id;
	`

	var summaryID string
	err := s.db.QueryRowContext(ctx, query, totalAmount, updatedAt, sessionId).Scan(&summaryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("no session found to update")
		}
		return "", err
	}

	return summaryID, nil
}

func (s *postgresRepository) CloseSession(ctx context.Context, sessionId string, sessionClosed bool, updatedAt time.Time) error {
	query := `UPDATE session_summary SET session_closed = $1, updated_at = $2
				WHERE session_id = $3;
			`
	res, err := s.db.ExecContext(ctx, query, sessionClosed, updatedAt, sessionId)
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

func (s *postgresRepository) AddOrder(ctx context.Context, sessionId string, newItems []pkg.OrderedItem) (pkg.Summary, error) {
	itemsJSON, err := json.Marshal(newItems)
	if err != nil {
		return pkg.Summary{}, fmt.Errorf("failed to marshal ordered items: %w", err)
	}

	var totalAmount float32
	totalAmount = 0
	log.Println("Adding the following orders")
	for _, items := range newItems {
		log.Println(items)
		totalAmount += items.ItemPrice * float32(items.Quantity)
	}

	var summary pkg.Summary
	var orderedItemsJSON []byte
	query := `
		UPDATE session_summary
		SET ordered_items = ordered_items || $2::jsonb,
			session_total_amount = session_total_amount + $3,
			updated_at = NOW()
		WHERE session_id = $1
		RETURNING summary_id,shop_id,session_id,session_number,session_total_amount,ordered_items,created_at,updated_at;
	`

	err = s.db.QueryRowContext(ctx, query, sessionId, itemsJSON, totalAmount).Scan(&summary.Id, &summary.ShopId, &summary.SessionId, &summary.SessionNumber, &summary.TotalAmount, &orderedItemsJSON, &summary.CreatedAt, &summary.UpdatedAt)
	if err != nil {
		return pkg.Summary{}, fmt.Errorf("failed to update ordered_items: %w", err)
	}

	if len(orderedItemsJSON) > 0 {
		var orderedItems []pkg.OrderedItem
		err = json.Unmarshal(orderedItemsJSON, &orderedItems)
		if err != nil {
			log.Printf("failed to unmarshal ordered_items for summary %s: %v", summary.Id, err)
			orderedItems = []pkg.OrderedItem{}
		}
		summary.OrderedItems = orderedItems
	}

	return summary, nil
}

func (s *postgresRepository) RemoveOrderItems(ctx context.Context, sessionId string, orderItemIds []string) error {
	orderItemIdsJSONB, err := json.Marshal(orderItemIds)
	if err != nil {
		return fmt.Errorf("failed to marshal orderItemIds: %w", err)
	}

	query := `
		UPDATE session_summary
		SET ordered_items = (
			SELECT jsonb_agg(item) 
			FROM jsonb_array_elements(ordered_items) AS item
			WHERE item->>'orderItemId' NOT IN (SELECT * FROM jsonb_array_elements($2::jsonb))
		)
		WHERE session_id = $1
	`
	res, err := s.db.ExecContext(ctx, query, sessionId, orderItemIdsJSONB)
	if err != nil {
		return fmt.Errorf("failed to execute remove order items query: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no session found or no rows updated")
	}

	return nil
}

func (s *postgresRepository) Close() error {
	return s.db.Close()
}
