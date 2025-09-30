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
	//MENU
	CreateMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetMenus(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	SetMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) error

	//ITEMS
	CreateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error

	//MENU-ITEMS
	CreateMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetMenuItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteMenuItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error

	//MENU-ITEMS
	CreateCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetCategories(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error
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

func (c *MenuClient) SetMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuId := vars["menuid"]

	res, err := c.service.SetMenu(ctx, &menu.SetMenuRequest{
		MenuId: menuId,
	})
	if err != nil {
		return err
	}

	return pkg.WriteJSON(w, http.StatusOK, pkg.SetMenuResponse{
		Success: res.Success,
	})
}

func (c *MenuClient) DeleteMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuId := vars["menuid"]

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

func (c *MenuClient) CreateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.CreateItemRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.CreateItem(ctx, &menu.CreateItemRequest{
		Name:        req.ItemName,
		ShopId:      req.ShopId,
		Description: req.Description,
		CategoryId:  req.CategoryId,
		Price:       req.Price,
		Image:       req.Image,
		IsAvailable: req.Available,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.CreateItemResponse{
		Item: pkg.Item{
			ItemId:      res.Item.ItemId,
			ShopId:      res.Item.ShopId,
			ItemName:    res.Item.Name,
			Description: res.Item.Description,
			Category:    res.Item.Category,
			Price:       res.Item.Price,
			Image:       res.Item.Image,
			Available:   res.Item.IsAvailable,
		},
	})
}

func (c *MenuClient) GetItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	itemId := vars["itemid"]

	res, err := c.service.GetItem(ctx, &menu.GetItemRequest{
		ItemId: itemId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetItemResponse{
		Item: pkg.Item{
			ItemId:      res.Item.ItemId,
			ShopId:      res.Item.ShopId,
			ItemName:    res.Item.Name,
			Description: res.Item.Description,
			Category:    res.Item.Category,
			Price:       res.Item.Price,
			Image:       res.Item.Image,
			Available:   res.Item.IsAvailable,
		},
	})
}

func (c *MenuClient) GetItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	shopId := vars["shopid"]

	res, err := c.service.GetItems(ctx, &menu.GetItemsRequest{
		ShopId: shopId,
	})
	if err != nil {
		return err
	}
	items := []pkg.Item{}
	for _, it := range res.Items {
		item := pkg.Item{
			ItemId:      it.ItemId,
			ShopId:      it.ShopId,
			ItemName:    it.Name,
			Description: it.Description,
			Category:    it.Category,
			Price:       it.Price,
			Image:       it.Image,
			Available:   it.IsAvailable,
		}
		items = append(items, item)
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetItemsResponse{
		Items: items,
	})
}

func (c *MenuClient) UpdateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	item := pkg.UpdateItemRequest{}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		return err
	}

	res, err := c.service.UpdateItem(ctx, &menu.UpdateItemRequest{
		ItemId:      item.ItemId,
		ShopId:      item.ShopId,
		Name:        item.ItemName,
		Description: item.Description,
		CategoryId:  item.CategoryId,
		Price:       item.Price,
		Image:       item.Image,
		IsAvailable: item.Available,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.UpdateItemResponse{
		Success: res.Success,
	})
}

func (c *MenuClient) DeleteItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	itemId := vars["itemid"]

	res, err := c.service.DeleteItem(ctx, &menu.DeleteItemRequest{
		ItemId: itemId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.DeleteItemResponse{
		Success: res.Success,
	})
}

func (c *MenuClient) CreateMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.CreateMenuItemRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.CreateMenuItem(ctx, &menu.CreateMenuItemRequest{
		MenuId: req.MenuId,
		ItemId: req.ItemId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.CreateMenuItemResponse{
		MenuItem: pkg.MenuItem{
			MenuItemId: res.MenuItem.MenuItemId,
			MenuId:     res.MenuItem.MenuId,
			ItemId:     res.MenuItem.ItemId,
		},
	})
}

func (c *MenuClient) GetMenuItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuId := vars["menuid"]

	res, err := c.service.GetMenuItems(ctx, &menu.GetMenuItemsRequest{
		MenuId: menuId,
	})
	if err != nil {
		return err
	}
	items := []pkg.Item{}
	for _, i := range res.Items {
		item := pkg.Item{
			ItemId:      i.ItemId,
			ShopId:      i.ShopId,
			ItemName:    i.Name,
			Description: i.Description,
			Category:    i.Category,
			Price:       i.Price,
			Image:       i.Image,
			Available:   i.IsAvailable,
		}
		items = append(items, item)
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetMenuItemsResponse{
		Items: items,
	})
}

func (c *MenuClient) DeleteMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.DeleteMenuItemRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.DeleteMenuItem(ctx, &menu.DeleteMenuItemRequest{
		MenuId: req.MenuId,
		ItemId: req.ItemId,
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
	menuId := vars["menuid"]

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

func (c *MenuClient) CreateCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.CreateCategoryRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.CreateCategory(ctx, &menu.CreateCategoryRequest{
		ShopId: req.ShopId,
		Name:   req.Name,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.CreateCategoryResponse{
		Category: pkg.Category{
			CategoryId: res.Category.Id,
			ShopId:     res.Category.ShopId,
			Name:       res.Category.Name,
		},
	})
}

func (c *MenuClient) GetCategories(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	shopId := vars["shopid"]

	res, err := c.service.GetCategory(ctx, &menu.GetCategoryRequest{
		ShopId: shopId,
	})
	if err != nil {
		return err
	}
	categories := []pkg.Category{}
	for _, c := range res.Categories {
		category := pkg.Category{
			CategoryId: c.Id,
			ShopId:     c.ShopId,
			Name:       c.Name,
		}
		categories = append(categories, category)
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetCategoriesResponse{
		Categories: categories,
	})
}

func (c *MenuClient) DeleteCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.DeleteCategoryRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.DeleteCategory(ctx, &menu.DeleteCategoryRequest{
		Id:     req.CategoryId,
		ShopId: req.ShopId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.DeleteCategoryResponse{
		Success: res.Success,
	})
}
