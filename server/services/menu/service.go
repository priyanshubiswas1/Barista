package menu

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/priyanshubiswas1/barista/menuservice/pkg"
)

type IService interface {
	CreateMenu(ctx context.Context, shopId, name string) (string, error)
	GetMenus(ctx context.Context, shopId string) ([]pkg.Menu, error)
	SetMenu(ctx context.Context, menuId string) (bool, error)
	DeleteMenu(ctx context.Context, menuId string) (bool, error)
	CreateItem(ctx context.Context, item pkg.ItemRequest) (pkg.Item, error)
	GetItem(ctx context.Context, itemId string) (*pkg.Item, error)
	GetItems(ctx context.Context, shopId string) ([]pkg.Item, error)
	GetItemsAdmin(ctx context.Context, shopId string) ([]pkg.Item, error)
	UpdateItem(ctx context.Context, item pkg.ItemRequest) (bool, error)
	UpdateItemAvailable(ctx context.Context, item pkg.UpdateItemAvailableRequest) (bool, error)
	DeleteItem(ctx context.Context, itemId string) (bool, error)
	CreateMenuItem(ctx context.Context, menuItem pkg.MenuItem) (pkg.MenuItem, error)
	GetMenuItems(ctx context.Context, menuId string) ([]pkg.Item, error)
	GetMenuItemsAdmin(ctx context.Context, menuId string) ([]pkg.Item, error)
	DeleteMenuItem(ctx context.Context, menuId, itemId string) (bool, error)
	DeleteMenuItems(ctx context.Context, menuId string) (bool, error)
	CreateCategory(ctx context.Context, category pkg.Category) (pkg.Category, error)
	UpdateCategory(ctx context.Context, category pkg.Category) (pkg.Category, error)
	GetCategories(ctx context.Context, shopId string) ([]pkg.Category, error)
	GetCategoriesWithItem(ctx context.Context, shopId string) ([]pkg.CategoryWithItems, error)
	DeleteCategory(ctx context.Context, id, shopId string) (bool, error)
	UploadImage(ctx context.Context, imageData []byte, filename string) (string, error)
}

type service struct {
	repository      IRepository
	cloudinaryService *pkg.CloudinaryService
}

func NewService(repository IRepository, cloudinaryService *pkg.CloudinaryService) IService {
	return &service{
		repository:       repository,
		cloudinaryService: cloudinaryService,
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

func (s *service) CreateItem(ctx context.Context, item pkg.ItemRequest) (pkg.Item, error) {
	categoryName, err := s.repository.GetCategoryName(ctx, item.CategoryId, item.ShopId)
	if err != nil {
		return pkg.Item{}, err
	}
	item.ItemId = uuid.New().String()
	if err := s.repository.CreateItem(ctx, item.ItemId, item.ShopId, item.Name, item.Description,
		item.CategoryId, item.Veg, item.Price, item.Image, item.IsAvailable, true, time.Now()); err != nil {
		return pkg.Item{}, err
	}
	return pkg.Item{
		ItemId:      item.ItemId,
		ShopId:      item.ShopId,
		Name:        item.Name,
		Description: item.Description,
		Category:    categoryName,
		Veg:         item.Veg,
		Price:       item.Price,
		Image:       item.Image,
		IsAvailable: item.IsAvailable,
	}, nil
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

func (s *service) GetItemsAdmin(ctx context.Context, shopId string) ([]pkg.Item, error) {
	if shopId == "" {
		return nil, errors.New("shop ID is required")
	}
	if _, err := uuid.Parse(shopId); err != nil {
		return nil, errors.New("invalid shop ID format")
	}
	return s.repository.GetItemsAdmin(ctx, shopId)
}

func (s *service) UpdateItem(ctx context.Context, item pkg.ItemRequest) (bool, error) {
	if item.ItemId == "" {
		return false, errors.New("item ID is required")
	}
	if _, err := uuid.Parse(item.ItemId); err != nil {
		return false, errors.New("invalid item ID format")
	}
	if _, err := s.repository.GetCategoryName(ctx, item.CategoryId, item.ShopId); err != nil {
		return false, err
	}
	return s.repository.UpdateItem(ctx, item.ItemId, item.ShopId, item.Name, item.Description, item.CategoryId, item.Veg, item.Price, item.Image, item.IsAvailable, time.Now())
}

func (s *service) UpdateItemAvailable(ctx context.Context, item pkg.UpdateItemAvailableRequest) (bool, error) {
	if item.ItemId == "" {
		return false, errors.New("item ID is required")
	}
	if _, err := uuid.Parse(item.ItemId); err != nil {
		return false, errors.New("invalid item ID format")
	}
	return s.repository.UpdateItemAvailable(ctx, item.ItemId, item.ShopId, item.IsAvailable, time.Now())
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

func (s *service) GetMenuItemsAdmin(ctx context.Context, menuId string) ([]pkg.Item, error) {
	if menuId == "" {
		return nil, errors.New("menu ID is required")
	}
	if _, err := uuid.Parse(menuId); err != nil {
		return nil, errors.New("invalid menu ID format")
	}
	return s.repository.GetMenuItemsAdmin(ctx, menuId)
}

func (s *service) DeleteMenuItem(ctx context.Context, menuId, itemId string) (bool, error) {
	return s.repository.DeleteMenuItem(ctx, menuId, itemId)
}

func (s *service) DeleteMenuItems(ctx context.Context, menuId string) (bool, error) {
	return s.repository.DeleteMenuItems(ctx, menuId)
}

func (s *service) CreateCategory(ctx context.Context, category pkg.Category) (pkg.Category, error) {
	category.Id = uuid.New().String()
	if err := s.repository.CreateCategory(ctx, category.Id, category.ShopId, category.Name); err != nil {
		return pkg.Category{}, nil
	}
	return pkg.Category{
		Id:     category.Id,
		ShopId: category.ShopId,
		Name:   category.Name,
	}, nil
}

func (s *service) UpdateCategory(ctx context.Context, category pkg.Category) (pkg.Category, error) {
	log.Println("NAME SERVICE: ", category.Name)
	if err := s.repository.UpdateCategory(ctx, category.Id, category.ShopId, category.Name); err != nil {
		return pkg.Category{}, nil
	}
	return pkg.Category{
		Id:     category.Id,
		ShopId: category.ShopId,
		Name:   category.Name,
	}, nil
}

func (s *service) GetCategories(ctx context.Context, shopId string) ([]pkg.Category, error) {
	return s.repository.GetCategories(ctx, shopId)
}
func (s *service) GetCategoriesWithItem(ctx context.Context, shopId string) ([]pkg.CategoryWithItems, error) {
	return s.repository.GetCategoriesWithItem(ctx, shopId)
}

func (s *service) DeleteCategory(ctx context.Context, id, shopId string) (bool, error) {
	return s.repository.DeleteCategory(ctx, id, shopId)
}

func (s *service) UploadImage(ctx context.Context, imageData []byte, filename string) (string, error) {
	if s.cloudinaryService == nil {
		return "", errors.New("cloudinary service not configured")
	}
	return s.cloudinaryService.UploadImage(ctx, imageData, filename)
}
