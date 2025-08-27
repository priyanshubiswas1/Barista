package shop

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/priyanshubiswas1/shopservice/pkg"
)

type IService interface {
	CreateShop(ctx context.Context, ownerId, shopName string) (string, error)
	GetByOwnerId(ctx context.Context, ownerId string) ([]*pkg.Shop, error)
	GetByShopId(ctx context.Context, shopId string) (*pkg.Shop, error)
	UpdateShop(ctx context.Context, ownerId string, shop *pkg.Shop) error
	DeleteShop(ctx context.Context, ownerId, shopId string) error
}

type service struct {
	repository IRepository
}

func NewService(repository IRepository) (IService, error) {
	return &service{
		repository: repository,
	}, nil
}

func (s *service) CreateShop(ctx context.Context, ownerId, shopName string) (string, error) {
	if ownerId == "" {
		return "", errors.New("owner ID is required")
	}
	if _, err := uuid.Parse(ownerId); err != nil {
		return "", errors.New("invalid owner ID format")
	}
	if shopName == "" {
		return "", errors.New("shop name is required")
	}
	id := uuid.New().String()
	if err := s.repository.Create(ctx, ownerId, id, shopName, true, time.Now()); err != nil {
		return "", err
	}
	return id, nil
}

func (s *service) GetByOwnerId(ctx context.Context, ownerId string) ([]*pkg.Shop, error) {
	if ownerId == "" {
		return nil, errors.New("owner ID is required")
	}
	if _, err := uuid.Parse(ownerId); err != nil {
		return nil, errors.New("invalid owner ID format")
	}
	return s.repository.GetByOwnerId(ctx, ownerId)
}

func (s *service) GetByShopId(ctx context.Context, shopId string) (*pkg.Shop, error) {
	if shopId == "" {
		return nil, errors.New("shop ID is required")
	}
	if _, err := uuid.Parse(shopId); err != nil {
		return nil, errors.New("invalid shop ID format")
	}
	return s.repository.GetByShopId(ctx, shopId)
}

func (s *service) UpdateShop(ctx context.Context, ownerId string, shop *pkg.Shop) error {
	if shop.ShopId == "" {
		return errors.New("shop ID is required")
	}
	if _, err := uuid.Parse(shop.ShopId); err != nil {
		return errors.New("invalid shop ID format")
	}
	return s.repository.Update(ctx, ownerId, shop.ShopId, shop.ShopName)
}

func (s *service) DeleteShop(ctx context.Context, ownerId, shopId string) error {
	if shopId == "" {
		return errors.New("shop ID is required")
	}
	if _, err := uuid.Parse(shopId); err != nil {
		return errors.New("invalid shop ID format")
	}
	return s.repository.Delete(ctx, ownerId, shopId)
}
