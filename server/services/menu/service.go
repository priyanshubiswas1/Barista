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
	DeleteMenus(ctx context.Context, menuId string) (bool, error)
	CreateMenuItem(ctx context.Context, menuItem pkg.MenuItem) (string, error)
	GetMenuItems(ctx context.Context, menuId string) ([]pkg.MenuItem, error)
	GetMenuItem(ctx context.Context, menuItemId string) (*pkg.MenuItem, error)
	UpdateMenuItem(ctx context.Context, menuItem pkg.MenuItem) (bool, error)
	DeleteMenuItem(ctx context.Context, menuItemID string) (bool, error)
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

func (s *service) DeleteMenus(ctx context.Context, menuId string) (bool, error) {
	if menuId == "" {
		return false, errors.New("menu ID is required")
	}
	if _, err := uuid.Parse(menuId); err != nil {
		return false, errors.New("invalid menu ID format")
	}
	return s.repository.DeleteMenu(ctx, menuId)
}

func (s *service) CreateMenuItem(ctx context.Context, menuItem pkg.MenuItem) (string, error) {
	menuItem.MenuItemId = uuid.New().String()
	if err := s.repository.CreateMenuItem(ctx, menuItem.MenuItemId, menuItem.MenuId, menuItem.Name, menuItem.Description,
		menuItem.Category, menuItem.Price, menuItem.Image, menuItem.IsAvailable, true, time.Now()); err != nil {
		return "", err
	}
	return menuItem.MenuItemId, nil
}

func (s *service) GetMenuItems(ctx context.Context, menuId string) ([]pkg.MenuItem, error) {
	if menuId == "" {
		return nil, errors.New("menu ID is required")
	}
	if _, err := uuid.Parse(menuId); err != nil {
		return nil, errors.New("invalid menu ID format")
	}
	return s.repository.GetMenuItems(ctx, menuId)
}

func (s *service) GetMenuItem(ctx context.Context, menuItemId string) (*pkg.MenuItem, error) {
	if menuItemId == "" {
		return nil, errors.New("menu ID is required")
	}
	if _, err := uuid.Parse(menuItemId); err != nil {
		return nil, errors.New("invalid menu ID format")
	}
	menuItem, err := s.repository.GetMenuItem(ctx, menuItemId)
	if err != nil {
		return nil, err
	}
	return &menuItem, nil
}

func (s *service) UpdateMenuItem(ctx context.Context, menuItem pkg.MenuItem) (bool, error) {
	if menuItem.MenuItemId == "" {
		return false, errors.New("menu item ID is required")
	}
	if _, err := uuid.Parse(menuItem.MenuItemId); err != nil {
		return false, errors.New("invalid menu item ID format")
	}
	if menuItem.MenuId == "" {
		return false, errors.New("menu ID is required")
	}
	if _, err := uuid.Parse(menuItem.MenuId); err != nil {
		return false, errors.New("invalid menu ID format")
	}
	return s.repository.UpdateMenuItem(ctx, menuItem.MenuItemId, menuItem.MenuId, menuItem.Name, menuItem.Description, menuItem.Category, menuItem.Price, menuItem.Image, menuItem.IsAvailable, true, time.Now())
}

func (s *service) DeleteMenuItem(ctx context.Context, menuItemID string) (bool, error) {
	return s.repository.DeleteMenuItem(ctx, menuItemID)
}

func (s *service) DeleteMenuItems(ctx context.Context, menuId string) (bool, error) {
	return s.repository.DeleteMenuItems(ctx, menuId)
}
