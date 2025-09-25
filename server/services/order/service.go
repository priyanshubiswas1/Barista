package order

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/priyanshubiswas1/barista/orderservice/pkg"
)

type IService interface {
	GetOrdersBySessionId(ctx context.Context, sessionId string) ([]pkg.OrderItems, error)
	GetOrderById(ctx context.Context, orderId string) ([]pkg.OrderItems, error)
	CreateOrder(ctx context.Context, order pkg.OrderReq) (pkg.CreateOrderRes, error)
	UpdateOrder(ctx context.Context, orderId, orderStatus string) bool
	CancelOrder(ctx context.Context, orderId string) bool
}

type service struct {
	repository IRepository
}

func NewService(repository IRepository) (IService, error) {
	return &service{
		repository: repository,
	}, nil
}

func (s *service) GetOrdersBySessionId(ctx context.Context, sessionId string) ([]pkg.OrderItems, error) {
	return s.repository.GetOrdersBySessionId(ctx, sessionId)
}

func (s *service) GetOrderById(ctx context.Context, orderId string) ([]pkg.OrderItems, error) {
	return s.repository.GetOrderById(ctx, orderId)
}

func (s *service) CreateOrder(ctx context.Context, order pkg.OrderReq) (pkg.CreateOrderRes, error) {
	id := uuid.New().String()
	return s.repository.CreateOrder(ctx, id, order.SessionId, order.SpecialNotes, time.Now(), order.OrderItems)
}

func (s *service) UpdateOrder(ctx context.Context, orderId, orderStatus string) bool {
	if err := s.repository.UpdateOrder(ctx, orderId, orderStatus, time.Now()); err != nil {
		return true
	} else {
		return false
	}
}

func (s *service) CancelOrder(ctx context.Context, orderId string) bool {
	if err := s.repository.CancelOrder(ctx, orderId, time.Now()); err != nil {
		return true
	} else {
		return false
	}
}

func (s *service) GetAllSessions(ctx context.Context, shopId string) error {
	if shopId == "" {
		return errors.New("shop ID is required")
	}
	if _, err := uuid.Parse(shopId); err != nil {
		return errors.New("invalid shop ID format")
	}
	return nil
}
