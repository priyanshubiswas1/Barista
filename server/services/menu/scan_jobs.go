package menu

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ScanJob struct {
	ShopID   string `json:"shopId"`
	ItemID   string `json:"itemId"`
	ImageURL string `json:"imageUrl"`
}

type ScanPublisher struct {
	conn       *amqp.Connection
	exchange   string
	routingKey string
}

func NewScanPublisher(conn *amqp.Connection, exchange, routingKey string) *ScanPublisher {
	if conn == nil || exchange == "" || routingKey == "" {
		return nil
	}
	return &ScanPublisher{conn: conn, exchange: exchange, routingKey: routingKey}
}

func (p *ScanPublisher) Enqueue(ctx context.Context, job ScanJob) error {
	if p == nil || p.conn == nil {
		return errors.New("scan publisher not configured")
	}
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	// Keep publishes bounded even when caller uses Background.
	pubCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return ch.PublishWithContext(
		pubCtx,
		p.exchange,
		p.routingKey,
		true,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
}

