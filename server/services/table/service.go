package table

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/priyanshubiswas1/barista/tableservice/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

type IService interface {
	CreateTable(ctx context.Context, shopId, tableNumber string) (pkg.Table, error)
	GetTable(ctx context.Context, tableId, shopId string) (pkg.Table, error)
	GetTables(ctx context.Context, shopId string) ([]pkg.Table, error)
	UpdateTable(ctx context.Context, table pkg.Table) error
	DeleteTable(ctx context.Context, tableId, shopId string) error
}

type service struct {
	repository   IRepository
	rmqPubConn   *amqp.Connection
	rmqConsConn  *amqp.Connection
	clientUrl    string
	exchange     string
	consumerName string
}

func NewService(repository IRepository, conn *amqp.Connection, consumeConn *amqp.Connection, clientUrl, exchange, consumerName string) IService {
	return &service{
		repository:   repository,
		rmqPubConn:   conn,
		rmqConsConn:  consumeConn,
		clientUrl:    clientUrl,
		exchange:     exchange,
		consumerName: consumerName,
	}
}

func (s *service) CreateTable(ctx context.Context, shopId, tableNumber string) (pkg.Table, error) {
	tableId := uuid.New().String()
	//CONSUME
	rmqConsumeClient, err := NewRabbitMQClient(s.rmqConsConn)
	if err != nil {
		return pkg.Table{}, err
	}
	defer rmqConsumeClient.Close()

	consumeQueue, err := rmqConsumeClient.CreateQueue("", true, true)
	if err != nil {
		return pkg.Table{}, err
	}

	if err := rmqConsumeClient.CreateBinding(consumeQueue.Name, consumeQueue.Name, "gen_qr_event"); err != nil {
		return pkg.Table{}, err
	}
	messageBus, err := rmqConsumeClient.StartConsumer(consumeQueue.Name, "table_service", true)
	if err != nil {
		return pkg.Table{}, err
	}
	correlationID := uuid.New().String()

	//PUBLISH
	rmqClient, err := NewRabbitMQClient(s.rmqPubConn)
	if err != nil {
		return pkg.Table{}, err
	}
	defer rmqClient.Close()

	if err := rmqClient.Send(ctx, "gen_qr_event", "gen_qr", amqp.Publishing{
		ContentType:   "text/plain",
		DeliveryMode:  amqp.Transient,
		ReplyTo:       consumeQueue.Name,
		CorrelationId: correlationID,
		Body:          fmt.Appendf([]byte{}, "%s/%s/%s", s.clientUrl, shopId, tableId),
	}); err != nil {
		return pkg.Table{}, err
	}

	var wg sync.WaitGroup
	var qr []byte

	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range messageBus {
			if msg.CorrelationId != correlationID {
				continue // not our message
			}
			if strings.HasPrefix(string(msg.Body), "error:") {
				log.Printf("Failed to generate QR: %s", msg.Body)
			} else {
				log.Printf("Generated QR %s", msg.Body)
				qr = msg.Body
			}
			if err := rmqConsumeClient.DeleteQueue(consumeQueue.Name, false, false); err != nil {
				log.Println("Failed to delete queue: ", err)
			}
			break
		}
	}()
	wg.Wait()

	if strings.HasPrefix(string(qr), "error:") {
		return pkg.Table{}, errors.New(string(qr))
	}
	encodedQRCode := base64.StdEncoding.EncodeToString(qr)
	createdAt := time.Now()

	if err := s.repository.CreateTable(ctx, tableId, shopId, tableNumber, encodedQRCode, true, createdAt); err != nil {
		return pkg.Table{}, err
	}
	return pkg.Table{
		TableId:     tableId,
		ShopId:      shopId,
		TableNumber: tableNumber,
		QR:          encodedQRCode,
	}, nil
}

func (s *service) GetTable(ctx context.Context, tableId, shopId string) (pkg.Table, error) {
	return s.repository.GetTable(ctx, tableId, shopId)
}

func (s *service) GetTables(ctx context.Context, shopId string) ([]pkg.Table, error) {
	return s.repository.GetTables(ctx, shopId)
}

func (s *service) UpdateTable(ctx context.Context, table pkg.Table) error {
	return s.repository.UpdateTable(ctx, table.TableId, table.ShopId, table.TableNumber, true, time.Now())
}

func (s *service) DeleteTable(ctx context.Context, tableId, shopId string) error {
	return s.repository.DeleteTable(ctx, tableId, shopId)
}
