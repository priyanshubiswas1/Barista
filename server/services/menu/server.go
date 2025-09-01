package menu

import (
	"context"
	"fmt"
	"net"

	"github.com/priyanshubiswas1/barista/menuservice/pkg"
	menuPb "github.com/priyanshubiswas1/baristacontracts/gen/menu"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCserver struct {
	service IService
	menuPb.UnimplementedMenuServiceServer
}

func ListenGRPC(service IService, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	menuPb.RegisterMenuServiceServer(server, &gRPCserver{service: service})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *gRPCserver) CreateMenu(ctx context.Context, req *menuPb.CreateMenuRequest) (*menuPb.CreateMenuResponse, error) {
	id, err := s.service.CreateMenu(ctx, req.ShopId, req.Name)
	if err != nil {
		return nil, err
	}
	return &menuPb.CreateMenuResponse{
		Menu: &menuPb.Menu{
			MenuId: id,
			ShopId: req.ShopId,
			Name:   req.Name,
		},
	}, nil
}

func (s *gRPCserver) GetMenus(ctx context.Context, req *menuPb.GetMenusRequest) (*menuPb.GetMenusResponse, error) {
	pkgMenus, err := s.service.GetMenus(ctx, req.ShopId)
	if err != nil {
		return nil, err
	}
	menus := []*menuPb.Menu{}
	for _, pkgMenu := range pkgMenus {
		menu := menuPb.Menu{
			MenuId: pkgMenu.MenuId,
			ShopId: pkgMenu.ShopId,
			Name:   pkgMenu.Name,
		}
		menus = append(menus, &menu)
	}
	return &menuPb.GetMenusResponse{
		Menus: menus,
	}, nil
}

func (s *gRPCserver) SetMenu(ctx context.Context, req *menuPb.SetMenuRequest) (*menuPb.SetMenuResponse, error) {
	res, err := s.service.SetMenu(ctx, req.MenuId)
	if err != nil {
		return nil, err
	}
	return &menuPb.SetMenuResponse{
		Success: res,
	}, nil
}

func (s *gRPCserver) DeleteMenu(ctx context.Context, req *menuPb.DeleteMenuRequest) (*menuPb.DeleteMenuResponse, error) {
	res, err := s.service.DeleteMenu(ctx, req.MenuId)
	if err != nil {
		return nil, err
	}
	return &menuPb.DeleteMenuResponse{
		Success: res,
	}, nil
}

func (s *gRPCserver) CreateItem(ctx context.Context, req *menuPb.CreateItemRequest) (*menuPb.CreateItemResponse, error) {
	item, err := s.service.CreateItem(ctx, pkg.ItemRequest{
		ItemId:      "",
		ShopId:      req.ShopId,
		Name:        req.Name,
		Description: req.Description,
		CategoryId:  req.CategoryId,
		Price:       req.Price,
		Image:       req.Image,
		IsAvailable: req.IsAvailable,
	})
	if err != nil {
		return nil, err
	}
	return &menuPb.CreateItemResponse{
		Item: &menuPb.Item{
			ItemId:      item.ItemId,
			ShopId:      item.ShopId,
			Name:        item.Name,
			Description: item.Description,
			Category:    item.Category,
			Price:       item.Price,
			Image:       item.Image,
			IsAvailable: item.IsAvailable,
		},
	}, nil
}

func (s *gRPCserver) GetItem(ctx context.Context, req *menuPb.GetItemRequest) (*menuPb.GetItemResponse, error) {
	item, err := s.service.GetItem(ctx, req.ItemId)
	if err != nil {
		return nil, err
	}
	return &menuPb.GetItemResponse{
		Item: &menuPb.Item{
			ItemId:      item.ItemId,
			ShopId:      item.ShopId,
			Name:        item.Name,
			Description: item.Description,
			Category:    item.Category,
			Price:       item.Price,
			Image:       item.Image,
			IsAvailable: item.IsAvailable,
		},
	}, nil
}

func (s *gRPCserver) GetItems(ctx context.Context, req *menuPb.GetItemsRequest) (*menuPb.GetItemsResponse, error) {
	pkgItems, err := s.service.GetItems(ctx, req.ShopId)
	if err != nil {
		return nil, err
	}
	items := []*menuPb.Item{}
	for _, item := range pkgItems {
		items = append(items, &menuPb.Item{
			ItemId:      item.ItemId,
			ShopId:      item.ShopId,
			Name:        item.Name,
			Description: item.Description,
			Category:    item.Category,
			Price:       item.Price,
			Image:       item.Image,
			IsAvailable: item.IsAvailable,
		})
	}
	return &menuPb.GetItemsResponse{
		Items: items,
	}, nil
}

func (s *gRPCserver) UpdateItem(ctx context.Context, req *menuPb.UpdateItemRequest) (*menuPb.UpdateItemResponse, error) {
	res, err := s.service.UpdateItem(ctx, pkg.ItemRequest{
		ItemId:      req.ItemId,
		ShopId:      req.ShopId,
		Name:        req.Name,
		Description: req.Description,
		CategoryId:  req.CategoryId,
		Price:       req.Price,
		Image:       req.Image,
		IsAvailable: req.IsAvailable,
	})
	if err != nil {
		return nil, err
	}
	return &menuPb.UpdateItemResponse{
		Success: res,
	}, nil
}

func (s *gRPCserver) DeleteItem(ctx context.Context, req *menuPb.DeleteItemRequest) (*menuPb.DeleteItemResponse, error) {
	res, err := s.service.DeleteItem(ctx, req.ItemId)
	if err != nil {
		return nil, err
	}
	return &menuPb.DeleteItemResponse{
		Success: res,
	}, nil
}

func (s *gRPCserver) CreateMenuItem(ctx context.Context, req *menuPb.CreateMenuItemRequest) (*menuPb.CreateMenuItemResponse, error) {
	menuItem, err := s.service.CreateMenuItem(ctx, pkg.MenuItem{
		Id:     "",
		MenuId: req.MenuId,
		ItemId: req.ItemId,
	})
	if err != nil {
		return nil, err
	}
	return &menuPb.CreateMenuItemResponse{
		MenuItem: &menuPb.MenuItem{
			MenuItemId: menuItem.Id,
			MenuId:     menuItem.MenuId,
			ItemId:     menuItem.ItemId,
		},
	}, nil
}

func (s *gRPCserver) GetMenuItems(ctx context.Context, req *menuPb.GetMenuItemsRequest) (*menuPb.GetMenuItemsResponse, error) {
	pkgItems, err := s.service.GetMenuItems(ctx, req.MenuId)
	if err != nil {
		return nil, err
	}
	items := []*menuPb.Item{}
	for _, pkgItem := range pkgItems {
		item := menuPb.Item{
			ItemId:      pkgItem.ItemId,
			ShopId:      pkgItem.ShopId,
			Name:        pkgItem.Name,
			Description: pkgItem.Description,
			Category:    pkgItem.Category,
			Price:       pkgItem.Price,
			Image:       pkgItem.Image,
			IsAvailable: pkgItem.IsAvailable,
		}
		items = append(items, &item)
	}
	return &menuPb.GetMenuItemsResponse{
		Items: items,
	}, nil
}

func (s *gRPCserver) DeleteMenuItem(ctx context.Context, req *menuPb.DeleteMenuItemRequest) (*menuPb.DeleteMenuItemResponse, error) {
	res, err := s.service.DeleteMenuItem(ctx, req.MenuId, req.ItemId)
	if err != nil {
		return nil, err
	}
	return &menuPb.DeleteMenuItemResponse{
		Success: res,
	}, nil
}

func (s *gRPCserver) DeleteMenuItems(ctx context.Context, req *menuPb.DeleteMenuItemsRequest) (*menuPb.DeleteMenuItemsResponse, error) {
	res, err := s.service.DeleteMenuItems(ctx, req.MenuId)
	if err != nil {
		return nil, err
	}
	return &menuPb.DeleteMenuItemsResponse{
		Success: res,
	}, nil
}

func (s *gRPCserver) CreateCategory(ctx context.Context, req *menuPb.CreateCategoryRequest) (*menuPb.CreateCategoryResponse, error) {
	category, err := s.service.CreateCategory(ctx, pkg.Category{
		ShopId: req.ShopId,
		Name:   req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &menuPb.CreateCategoryResponse{
		Category: &menuPb.Category{
			Id:     category.Id,
			ShopId: category.ShopId,
			Name:   category.Name,
		},
	}, nil
}

func (s *gRPCserver) GetCategory(ctx context.Context, req *menuPb.GetCategoryRequest) (*menuPb.GetCategoryResponse, error) {
	pkgCategories, err := s.service.GetCategories(ctx, req.ShopId)
	if err != nil {
		return nil, err
	}
	categories := []*menuPb.Category{}
	for _, c := range pkgCategories {
		categories = append(categories, &menuPb.Category{
			Id:     c.Id,
			ShopId: c.ShopId,
			Name:   c.Name,
		})
	}
	return &menuPb.GetCategoryResponse{
		Categories: categories,
	}, nil
}

func (s *gRPCserver) DeleteCategory(ctx context.Context, req *menuPb.DeleteCategoryRequest) (*menuPb.DeleteCategoryResponse, error) {
	res, err := s.service.DeleteCategory(ctx, req.Id, req.ShopId)
	if err != nil {
		return nil, err
	}
	return &menuPb.DeleteCategoryResponse{
		Success: res,
	}, nil
}
