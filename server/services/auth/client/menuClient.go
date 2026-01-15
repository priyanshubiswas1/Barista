package client

import (
	"context"

	"errors"

	"github.com/priyanshubiswas1/barista/authservice/pkg"
	menu "github.com/priyanshubiswas1/baristacontracts/gen/menu"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IMenuClient interface {
	CreateMenu(ctx context.Context, menuReq pkg.CreateMenuReq) (string, error)
	GetMenus(ctx context.Context, shopId string) ([]string, error)
}

type MenuClient struct {
	client  *grpc.ClientConn
	service menu.MenuServiceClient
}

func NewMenu(url string) (IMenuClient, error) {
	if url == "" {
		return &MenuClient{}, nil
	}
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &MenuClient{
		client:  client,
		service: menu.NewMenuServiceClient(client),
	}, nil
}

func (c *MenuClient) CreateMenu(ctx context.Context, menuReq pkg.CreateMenuReq) (string, error) {
	if c.service == nil {
		return "", errors.New("menu service URL not configured")
	}
	res, err := c.service.CreateMenu(ctx, &menu.CreateMenuRequest{
		ShopId: menuReq.ShopId,
		Name:   menuReq.MenuName,
	})
	if err != nil {
		return "", err
	}
	return res.Menu.MenuId, err
}

func (c *MenuClient) GetMenus(ctx context.Context, shopId string) ([]string, error) {
	if c.service == nil {
		return nil, errors.New("menu service URL not configured")
	}
	res, err := c.service.GetMenus(ctx, &menu.GetMenusRequest{
		ShopId: shopId,
	})
	if err != nil {
		return nil, err
	}
	menuIds := []string{}
	for _, menu := range res.Menus {
		menuIds = append(menuIds, menu.MenuId)
	}
	return menuIds, nil
}
