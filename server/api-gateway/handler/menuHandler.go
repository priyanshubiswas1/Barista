package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/priyanshubiswas1/barista/apigateway/pkg"
	menu "github.com/priyanshubiswas1/baristacontracts/gen/menu"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IMenuClient interface {
	CreateMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetMenus(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	CreateMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetMenuItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteMenuItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

type MenuClient struct {
	client  *grpc.ClientConn
	service menu.MenuServiceClient
}

func NewMenu(url string) (IMenuClient, error) {
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &MenuClient{
		client:  client,
		service: menu.NewMenuServiceClient(client),
	}, nil
}

func (c *MenuClient) CreateMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.CreateMenuRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.CreateMenu(ctx, &menu.CreateMenuRequest{
		ShopId: req.ShopId,
		Name:   req.Name,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.CreateMenuResponse{
		Menu: pkg.Menu{
			MenuId: res.Menu.MenuId,
			ShopId: res.Menu.ShopId,
			Name:   res.Menu.Name,
		},
	})
}

func (c *MenuClient) GetMenus(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	shopId := vars["shopid"]

	res, err := c.service.GetMenus(ctx, &menu.GetMenusRequest{
		ShopId: shopId,
	})
	if err != nil {
		return err
	}

	menus := []pkg.Menu{}
	for _, m := range res.Menus {
		menu := pkg.Menu{
			MenuId: m.MenuId,
			ShopId: m.ShopId,
			Name:   m.Name,
		}
		menus = append(menus, menu)
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetMenusResponse{
		Menus: menus,
	})
}

func (c *MenuClient) DeleteMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuId := vars["menuId"]

	res, err := c.service.DeleteMenu(ctx, &menu.DeleteMenuRequest{
		MenuId: menuId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.DeleteMenuResponse{
		Success: res.Success,
	})
}

func (c *MenuClient) CreateMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuId := vars["menuId"]

	req := pkg.CreateMenuItemRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.CreateMenuItem(ctx, &menu.CreateMenuItemRequest{
		MenuId:      menuId,
		Name:        req.ItemName,
		Description: req.Description,
		Category:    req.Category,
		Price:       req.Price,
		Image:       req.Image,
		IsAvailable: req.Available,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.CreateMenuItemResponse{
		MenuItem: pkg.MenuItem{
			MenuItemId:  res.MenuItem.MenuItemId,
			MenuId:      res.MenuItem.MenuId,
			ItemName:    res.MenuItem.Name,
			Description: res.MenuItem.Description,
			Category:    res.MenuItem.Category,
			Price:       res.MenuItem.Price,
			Image:       res.MenuItem.Image,
			Available:   res.MenuItem.IsAvailable,
		},
	})
}

func (c *MenuClient) GetMenuItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuId := vars["menuId"]

	res, err := c.service.GetMenuItems(ctx, &menu.GetMenuItemsRequest{
		MenuId: menuId,
	})
	if err != nil {
		return err
	}
	menuItems := []pkg.MenuItem{}
	for _, i := range res.MenuItems {
		menuItem := pkg.MenuItem{
			MenuItemId:  i.MenuItemId,
			MenuId:      i.MenuId,
			ItemName:    i.Name,
			Description: i.Description,
			Category:    i.Category,
			Price:       i.Price,
			Image:       i.Image,
			Available:   i.IsAvailable,
		}
		menuItems = append(menuItems, menuItem)
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetMenuItemsResponse{
		MenuItems: menuItems,
	})
}

func (c *MenuClient) GetMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuItemId := vars["menu-itemId"]

	res, err := c.service.GetMenuItem(ctx, &menu.GetMenuItemRequest{
		MenuItemId: menuItemId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetMenuItemResponse{
		MenuItem: pkg.MenuItem{
			MenuItemId:  res.MenuItem.MenuItemId,
			MenuId:      res.MenuItem.MenuId,
			ItemName:    res.MenuItem.Name,
			Description: res.MenuItem.Description,
			Category:    res.MenuItem.Category,
			Price:       res.MenuItem.Price,
			Image:       res.MenuItem.Image,
			Available:   res.MenuItem.IsAvailable,
		},
	})
}

func (c *MenuClient) UpdateMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuItemId := vars["menu-itemId"]

	menuItem := pkg.MenuItem{}
	if err := json.NewDecoder(r.Body).Decode(&menuItem); err != nil {
		return err
	}

	menuItem.MenuItemId = menuItemId

	res, err := c.service.UpdateMenuItem(ctx, &menu.UpdateMenuItemRequest{
		MenuItem: &menu.MenuItem{
			MenuItemId:  menuItem.MenuItemId,
			MenuId:      menuItem.MenuId,
			Name:        menuItem.ItemName,
			Description: menuItem.Description,
			Category:    menuItem.Category,
			Price:       menuItem.Price,
			Image:       menuItem.Image,
			IsAvailable: menuItem.Available,
		},
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.UpdateMenuItemResponse{
		Success: res.Success,
	})
}

func (c *MenuClient) DeleteMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuItemId := vars["menu-itemId"]

	res, err := c.service.DeleteMenuItem(ctx, &menu.DeleteMenuItemRequest{
		MenuItemId: menuItemId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.DeleteMenuItemResponse{
		Success: res.Success,
	})
}

func (c *MenuClient) DeleteMenuItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuId := vars["menuId"]

	res, err := c.service.DeleteMenuItems(ctx, &menu.DeleteMenuItemsRequest{
		MenuId: menuId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.DeleteMenuItemsResponse{
		Success: res.Success,
	})
}
