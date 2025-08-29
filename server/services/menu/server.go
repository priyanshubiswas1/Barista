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

func (s *gRPCserver) DeleteMenu(ctx context.Context, req *menuPb.DeleteMenuRequest) (*menuPb.DeleteMenuResponse, error) {
	res, err := s.service.DeleteMenus(ctx, req.MenuId)
	if err != nil {
		return nil, err
	}
	return &menuPb.DeleteMenuResponse{
		Success: res,
	}, nil
}

func (s *gRPCserver) CreateMenuItem(ctx context.Context, req *menuPb.CreateMenuItemRequest) (*menuPb.CreateMenuItemResponse, error) {
	id, err := s.service.CreateMenuItem(ctx, pkg.MenuItem{
		MenuId:      req.MenuId,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Price:       req.Price,
		Image:       req.Image,
		IsAvailable: req.IsAvailable,
	})
	if err != nil {
		return nil, err
	}
	return &menuPb.CreateMenuItemResponse{
		MenuItem: &menuPb.MenuItem{
			MenuItemId:  id,
			MenuId:      req.MenuId,
			Name:        req.Name,
			Description: req.Description,
			Category:    req.Category,
			Price:       req.Price,
			Image:       req.Image,
			IsAvailable: req.IsAvailable,
		},
	}, nil
}

func (s *gRPCserver) GetMenuItems(ctx context.Context, req *menuPb.GetMenuItemsRequest) (*menuPb.GetMenuItemsResponse, error) {
	pkgMenuItems, err := s.service.GetMenuItems(ctx, req.MenuId)
	if err != nil {
		return nil, err
	}
	menuItems := []*menuPb.MenuItem{}
	for _, pkgMenuItem := range pkgMenuItems {
		menuItem := menuPb.MenuItem{
			MenuItemId:  pkgMenuItem.MenuItemId,
			MenuId:      pkgMenuItem.MenuId,
			Name:        pkgMenuItem.Name,
			Description: pkgMenuItem.Description,
			Category:    pkgMenuItem.Category,
			Price:       pkgMenuItem.Price,
			Image:       pkgMenuItem.Image,
			IsAvailable: pkgMenuItem.IsAvailable,
		}
		menuItems = append(menuItems, &menuItem)
	}
	return &menuPb.GetMenuItemsResponse{
		MenuItems: menuItems,
	}, nil
}

func (s *gRPCserver) GetMenuItem(ctx context.Context, req *menuPb.GetMenuItemRequest) (*menuPb.GetMenuItemResponse, error) {
	pkgMenuItem, err := s.service.GetMenuItem(ctx, req.MenuItemId)
	if err != nil {
		return nil, err
	}
	return &menuPb.GetMenuItemResponse{
		MenuItem: &menuPb.MenuItem{
			MenuItemId:  pkgMenuItem.MenuItemId,
			MenuId:      pkgMenuItem.MenuId,
			Name:        pkgMenuItem.Name,
			Description: pkgMenuItem.Description,
			Category:    pkgMenuItem.Category,
			Price:       pkgMenuItem.Price,
			Image:       pkgMenuItem.Image,
			IsAvailable: pkgMenuItem.IsAvailable,
		},
	}, nil
}

func (s *gRPCserver) UpdateMenuItem(ctx context.Context, req *menuPb.UpdateMenuItemRequest) (*menuPb.UpdateMenuItemResponse, error) {
	res, err := s.service.UpdateMenuItem(ctx, pkg.MenuItem{
		MenuItemId:  req.MenuItem.MenuItemId,
		MenuId:      req.MenuItem.MenuId,
		Name:        req.MenuItem.Name,
		Description: req.MenuItem.Description,
		Category:    req.MenuItem.Category,
		Price:       req.MenuItem.Price,
		Image:       req.MenuItem.Image,
		IsAvailable: req.MenuItem.IsAvailable,
	})
	if err != nil {
		return nil, err
	}
	return &menuPb.UpdateMenuItemResponse{
		Success: res,
	}, nil
}

func (s *gRPCserver) DeleteMenuItem(ctx context.Context, req *menuPb.DeleteMenuItemRequest) (*menuPb.DeleteMenuItemResponse, error) {
	res, err := s.service.DeleteMenuItem(ctx, req.MenuItemId)
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
