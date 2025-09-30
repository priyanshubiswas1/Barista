package payment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/priyanshubiswas1/barista/paymentservice/client"
	"github.com/priyanshubiswas1/barista/paymentservice/pkg"
)

type IService interface {
	GetPayment(ctx context.Context, sessionId string) (pkg.Payment, error)
	CreatePayment(ctx context.Context, sessionId, paymentMethod, paymentStatus, transactionId string, amount float32) (string, error)
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

func (s *service) GetPayment(ctx context.Context, sessionId string) (pkg.Payment, error) {
	if sessionId == "" {
		return pkg.Payment{}, errors.New("session ID is required")
	}
	if _, err := uuid.Parse(sessionId); err != nil {
		return pkg.Payment{}, errors.New("invalid session ID format")
	}
	return s.repository.GetPayment(ctx, sessionId)
}

func (s *service) CreatePayment(ctx context.Context, sessionId, paymentMethod, paymentStatus, transactionId string, amount float32) (string, error) {
	id := uuid.New().String()
	session, err := s.client.GetSessionById(ctx, sessionId)
	if err != nil {
		return "", err
	}

	if amount != session.Amount {
		return "", fmt.Errorf("paid amount mismatch: expected %.2f, got %.2f", session.Amount, amount)
	}

	if err := s.repository.CreatePayment(ctx, id, sessionId, paymentMethod, paymentStatus, transactionId, amount, time.Now()); err != nil {
		return "", err
	}

	succ, err := s.client.UpdateSession(ctx, pkg.UpdateSessionRequest{
		Session: pkg.Session{
			Id:            sessionId,
			ShopId:        "",
			TableId:       "",
			Otp:           "",
			Status:        "",
			Amount:        0.0,
			IsPaid:        true,
			PaymentMethod: paymentMethod,
			PhoneNumber:   "",
			Active:        true,
		},
	})
	if err != nil {
		return "", err
	}
	if !succ {
		return "", fmt.Errorf("failed to update session status")
	}
	return id, nil
}
