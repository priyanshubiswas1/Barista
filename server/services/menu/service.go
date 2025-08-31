package menu

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/priyanshubiswas1/barista/menuservice/pkg"
)

type IService interface {
	CreateMenu(ctx context.Context, shopId, name string) (string, error)
	GetMenus(ctx context.Context, shopId string) ([]pkg.Menu, error)
	SetMenu(ctx context.Context, menuId string) (bool, error)
	DeleteMenu(ctx context.Context, menuId string) (bool, error)
	CreateItem(ctx context.Context, item pkg.Item) (string, error)
	GetItem(ctx context.Context, itemId string) (*pkg.Item, error)
	GetItems(ctx context.Context, shopId string) ([]pkg.Item, error)
	UpdateItem(ctx context.Context, item pkg.Item) (bool, error)
	DeleteItem(ctx context.Context, itemId string) (bool, error)
	CreateMenuItem(ctx context.Context, menuItem pkg.MenuItem) (pkg.MenuItem, error)
	GetMenuItems(ctx context.Context, menuId string) ([]pkg.Item, error)
	DeleteMenuItem(ctx context.Context, menuId, itemId string) (bool, error)
	DeleteMenuItems(ctx context.Context, menuId string) (bool, error)
}

type service struct {
	repository IRepository
}

func NewService(repository IRepository) IService {
	return &service{
		repository: repository,
	}
}

func (s *service) CreateMenu(ctx context.Context, shopId, name string) (string, error) {
	id := uuid.New().String()
	if err := s.repository.CreateMenu(ctx, id, shopId, name, true, time.Now()); err != nil {
		return "", err
	}
	return id, nil
}

func (s *service) GetMenus(ctx context.Context, shopId string) ([]pkg.Menu, error) {
	if shopId == "" {
		return nil, errors.New("shop ID is required")
	}
	if _, err := uuid.Parse(shopId); err != nil {
		return nil, errors.New("invalid shop ID format")
	}
	return s.repository.GetMenus(ctx, shopId)
}

func (s *service) SetMenu(ctx context.Context, menuId string) (bool, error) {
	if menuId == "" {
		return false, errors.New("menu ID is required")
	}
	if _, err := uuid.Parse(menuId); err != nil {
		return false, errors.New("invalid shop ID format")
	}
	return s.repository.SetMenu(ctx, menuId)
}

func (s *service) DeleteMenu(ctx context.Context, menuId string) (bool, error) {
	if menuId == "" {
		return false, errors.New("menu ID is required")
	}
	if _, err := uuid.Parse(menuId); err != nil {
		return false, errors.New("invalid menu ID format")
	}
	return s.repository.DeleteMenu(ctx, menuId)
}

func (s *service) CreateItem(ctx context.Context, item pkg.Item) (string, error) {
	item.ItemId = uuid.New().String()
	if err := s.repository.CreateItem(ctx, item.ItemId, item.ShopId, item.Name, item.Description,
		item.Category, item.Price, item.Image, item.IsAvailable, true, time.Now()); err != nil {
		return "", err
	}
	return item.ItemId, nil
}

func (s *service) GetItem(ctx context.Context, itemId string) (*pkg.Item, error) {
	if itemId == "" {
		return nil, errors.New("item ID is required")
	}
	if _, err := uuid.Parse(itemId); err != nil {
		return nil, errors.New("invalid item ID format")
	}
	item, err := s.repository.GetItem(ctx, itemId)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *service) GetItems(ctx context.Context, shopId string) ([]pkg.Item, error) {
	if shopId == "" {
		return nil, errors.New("shop ID is required")
	}
	if _, err := uuid.Parse(shopId); err != nil {
		return nil, errors.New("invalid shop ID format")
	}
	return s.repository.GetItems(ctx, shopId)
}

func (s *service) UpdateItem(ctx context.Context, item pkg.Item) (bool, error) {
	if item.ItemId == "" {
		return false, errors.New("item ID is required")
	}
	if _, err := uuid.Parse(item.ItemId); err != nil {
		return false, errors.New("invalid item ID format")
	}
	return s.repository.UpdateItem(ctx, item.ItemId, item.ShopId, item.Name, item.Description, item.Category, item.Price, item.Image, item.IsAvailable, true, time.Now())
}

func (s *service) DeleteItem(ctx context.Context, itemId string) (bool, error) {
	return s.repository.DeleteItem(ctx, itemId)
}

func (s *service) CreateMenuItem(ctx context.Context, menuItem pkg.MenuItem) (pkg.MenuItem, error) {
	menuItem.Id = uuid.New().String()
	menuItem, err := s.repository.CreateMenuItem(ctx, menuItem.Id, menuItem.MenuId, menuItem.ItemId)
	if err != nil {
		return pkg.MenuItem{}, err
	}
	return menuItem, nil
}

func (s *service) GetMenuItems(ctx context.Context, menuId string) ([]pkg.Item, error) {
	if menuId == "" {
		return nil, errors.New("menu ID is required")
	}
	if _, err := uuid.Parse(menuId); err != nil {
		return nil, errors.New("invalid menu ID format")
	}
	return s.repository.GetMenuItems(ctx, menuId)
}

func (s *service) DeleteMenuItem(ctx context.Context, menuId, itemId string) (bool, error) {
	return s.repository.DeleteMenuItem(ctx, menuId, itemId)
}

func (s *service) DeleteMenuItems(ctx context.Context, menuId string) (bool, error) {
	return s.repository.DeleteMenuItems(ctx, menuId)
}
