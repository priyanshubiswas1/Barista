package history

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/priyanshubiswas1/barista/historyservice/pkg"
)

type IService interface {
	GetSessionSummary(ctx context.Context, shopId string) ([]pkg.Summary, error)
	CreateSession(ctx context.Context, shopId, sessionId string, sessionNumber int32, totalAmount float32, orderedItems []byte) (string, time.Time, error)
	UpdateSession(ctx context.Context, sessionId string, totalAmount float32) (string, time.Time, error)
	CloseSession(ctx context.Context, sessionId string) error
	AddOrder(ctx context.Context, sessionId string, newItems []pkg.OrderedItem) (pkg.Summary, error)
	RemoveOrderItems(ctx context.Context, sessionId string, orderItemIds []string) error
}

type service struct {
	repository IRepository
}

func NewService(repository IRepository) (IService, error) {
	return &service{
		repository: repository,
	}, nil
}

func (s *service) GetSessionSummary(ctx context.Context, shopId string) ([]pkg.Summary, error) {
	if shopId == "" {
		return nil, errors.New("shop ID is required")
	}
	if _, err := uuid.Parse(shopId); err != nil {
		return nil, errors.New("invalid shop ID format")
	}
	return s.repository.GetSessionSummary(ctx, shopId)
}

func (s *service) CreateSession(ctx context.Context, shopId, sessionId string, sessionNumber int32, totalAmount float32, orderedItems []byte) (string, time.Time, error) {
	log.Println("HistoryService: CreateSession called")
	id := uuid.New().String()
	if shopId == "" {
		return "", time.Time{}, errors.New("shop ID is required")
	}
	if _, err := uuid.Parse(shopId); err != nil {
		return "", time.Time{}, errors.New("invalid shop ID format")
	}
	if sessionId == "" {
		return "", time.Time{}, errors.New("session ID is required")
	}
	if _, err := uuid.Parse(sessionId); err != nil {
		return "", time.Time{}, errors.New("invalid session ID format")
	}
	createdAt := time.Now()
	if err := s.repository.CreateSession(ctx, id, shopId, sessionId, sessionNumber, totalAmount, orderedItems, false, createdAt); err != nil {
		return "", time.Time{}, err
	}
	return id, createdAt, nil
}

func (s *service) CloseSession(ctx context.Context, sessionId string) error {
	if sessionId == "" {
		return errors.New("session ID is required")
	}
	if _, err := uuid.Parse(sessionId); err != nil {
		return errors.New("invalid session ID format")
	}
	return s.repository.CloseSession(ctx, sessionId, true, time.Now())
}

func (s *service) UpdateSession(ctx context.Context, sessionId string, totalAmount float32) (string, time.Time, error) {
	if sessionId == "" {
		return "", time.Time{}, errors.New("session ID is required")
	}
	if _, err := uuid.Parse(sessionId); err != nil {
		return "", time.Time{}, errors.New("invalid session ID format")
	}
	updatedAt := time.Now()
	summaryId, err := s.repository.UpdateSession(ctx, sessionId, totalAmount, updatedAt)
	if err != nil {
		return "", time.Time{}, err
	}
	return summaryId, updatedAt, nil
}

func (s *service) AddOrder(ctx context.Context, sessionId string, newItems []pkg.OrderedItem) (pkg.Summary, error) {
	if sessionId == "" {
		return pkg.Summary{}, errors.New("session ID is required")
	}
	if _, err := uuid.Parse(sessionId); err != nil {
		return pkg.Summary{}, errors.New("invalid session ID format")
	}
	return s.repository.AddOrder(ctx, sessionId, newItems)
}

func (s *service) RemoveOrderItems(ctx context.Context, sessionId string, orderItemIds []string) error {
	if sessionId == "" {
		return errors.New("session ID is required")
	}
	if _, err := uuid.Parse(sessionId); err != nil {
		return errors.New("invalid session ID format")
	}
	return s.repository.RemoveOrderItems(ctx, sessionId, orderItemIds)
}
