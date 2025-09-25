package payment

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/priyanshubiswas1/barista/paymentservice/pkg"
)

type IService interface {
	GetPayment(ctx context.Context, sessionId string) (pkg.Payment, error)
	CreatePayment(ctx context.Context, sessionId, paymentMethod, paymentStatus, transactionId string, amount float32) (string, error)
}

type service struct {
	repository IRepository
}

func NewService(repository IRepository) (IService, error) {
	return &service{
		repository: repository,
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
	if err := s.repository.CreatePayment(ctx, id, sessionId, paymentMethod, paymentStatus, transactionId, amount, time.Now()); err != nil {
		return "", err
	}
	return id, nil
}
