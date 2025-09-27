package order

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/priyanshubiswas1/barista/orderservice/client"
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
	client     client.ISessionClient
}

func NewService(repository IRepository, client client.ISessionClient) (IService, error) {
	return &service{
		repository: repository,
		client:     client,
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
	res, err := s.repository.CreateOrder(ctx, id, order.SessionId, order.SpecialNotes, time.Now(), order.OrderItems)
	if err != nil {
		return pkg.CreateOrderRes{}, err
	}
	var orderAmount float32 = 0.00
	for _, oi := range order.OrderItems {
		orderAmount += oi.PriceAtOrderTime * float32(oi.Quantity)
	}
	succ, err := s.client.UpdateSession(ctx, pkg.UpdateSessionRequest{
		Session: pkg.Session{
			Id:            order.SessionId,
			ShopId:        "",
			TableId:       "",
			Otp:           "",
			Status:        "",
			Amount:        orderAmount,
			IsPaid:        false,
			PaymentMethod: "",
			PhoneNumber:   "",
			Active:        true,
		},
	})
	if err != nil {
		return pkg.CreateOrderRes{}, err
	}
	if succ == true {
		return res, err
	}
	return res, errors.New("failed to update session table")
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
