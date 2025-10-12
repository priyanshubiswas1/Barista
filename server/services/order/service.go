package order

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/priyanshubiswas1/barista/orderservice/client"
	"github.com/priyanshubiswas1/barista/orderservice/pkg"
)

type IService interface {
	GetOrdersBySessionId(ctx context.Context, sessionId string) ([]pkg.OrderItems, float32, error)
	GetOrderById(ctx context.Context, orderId string) ([]pkg.OrderItems, float32, error)
	CreateOrder(ctx context.Context, order pkg.OrderReq) (pkg.CreateOrderRes, error)
	CreateOrderAdmin(ctx context.Context, order pkg.OrderAdminReq) (pkg.CreateOrderRes, error)
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

func (s *service) GetOrdersBySessionId(ctx context.Context, sessionId string) ([]pkg.OrderItems, float32, error) {
	orderItems, err := s.repository.GetOrdersBySessionId(ctx, sessionId)
	if err != nil {
		return nil, 0, err
	}
	var orderAmount float32 = 0.00
	for _, oi := range orderItems {
		if oi.OrderStatus != "cancelled" {
			orderAmount += oi.PriceAtOrderTime * float32(oi.Quantity)
		}
	}
	return orderItems, orderAmount, nil
}

func (s *service) GetOrderById(ctx context.Context, orderId string) ([]pkg.OrderItems, float32, error) {
	orderItems, err := s.repository.GetOrderById(ctx, orderId)
	if err != nil {
		return nil, 0, err
	}
	var orderAmount float32 = 0.00
	for _, oi := range orderItems {
		orderAmount += oi.PriceAtOrderTime * float32(oi.Quantity)
	}
	return orderItems, orderAmount, nil
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
		log.Println("failed to UpdateSession error: ", err)
		return pkg.CreateOrderRes{}, err
	}
	if succ {
		return res, err
	}
	return res, errors.New("failed to update session table")
}

func (s *service) CreateOrderAdmin(ctx context.Context, order pkg.OrderAdminReq) (pkg.CreateOrderRes, error) {
	session, _, err := s.client.CreateSessionAdmin(ctx, order.ShopId, order.TableId, "active", "", "0000000000", 0, false, true)
	if err != nil {
		log.Println("failed to Create or Join Session error: ", err)
		return pkg.CreateOrderRes{}, err
	}

	id := uuid.New().String()
	res, err := s.repository.CreateOrder(ctx, id, session.Id, order.SpecialNotes, time.Now(), order.OrderItems)
	if err != nil {
		return pkg.CreateOrderRes{}, err
	}
	var orderAmount float32 = 0.00
	for _, oi := range order.OrderItems {
		orderAmount += oi.PriceAtOrderTime * float32(oi.Quantity)
	}
	succ, err := s.client.UpdateSession(ctx, pkg.UpdateSessionRequest{
		Session: pkg.Session{
			Id:            session.Id,
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
		log.Println("failed to UpdateSession error: ", err)
		return pkg.CreateOrderRes{}, err
	}
	if succ {
		return res, err
	}
	return res, errors.New("failed to update session table")
}

func (s *service) UpdateOrder(ctx context.Context, orderId, orderStatus string) bool {
	if err := s.repository.UpdateOrder(ctx, orderId, orderStatus, time.Now()); err != nil {
		log.Println("failed to UpdateOrder error: ", err)
		return false
	} else {
		return true
	}
}

func (s *service) CancelOrder(ctx context.Context, orderId string) bool {
	_, cancelAmount, err := s.GetOrderById(ctx, orderId)
	if err != nil {
		log.Println("failed to GetOrderById error: ", err)
		return false
	}

	sessionId, err := s.repository.GetSessionIdByOrderId(ctx, orderId)
	if err != nil {
		log.Println("failed to get sessionIdByOrderId error: ", err)
		return false
	}

	if err := s.repository.CancelOrder(ctx, orderId, time.Now()); err != nil {
		log.Println("failed to CancelOrder error: ", err)
		return false
	}
	succ, err := s.client.UpdateSession(ctx, pkg.UpdateSessionRequest{
		Session: pkg.Session{
			Id:            sessionId,
			ShopId:        "",
			TableId:       "",
			Otp:           "",
			Status:        "",
			Amount:        -cancelAmount,
			IsPaid:        false,
			PaymentMethod: "",
			PhoneNumber:   "",
			Active:        true,
		},
	})
	if err != nil {
		log.Println("failed to UpdateSession error: ", err)
		return false
	}
	return succ
}
