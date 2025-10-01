package client

import (
	"context"

	"github.com/priyanshubiswas1/barista/authservice/pkg"
	shop "github.com/priyanshubiswas1/baristacontracts/gen/shop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IShopClient interface {
	RegisterShop(ctx context.Context, shop pkg.CreateShopReq) (string, error)
	GetShopByOwnerId(ctx context.Context, id string) ([]string, error)
}

type ShopClient struct {
	client  *grpc.ClientConn
	service shop.ShopServiceClient
}

func NewShop(url string) (IShopClient, error) {
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &ShopClient{
		client:  client,
		service: shop.NewShopServiceClient(client),
	}, nil
}

func (c *ShopClient) RegisterShop(ctx context.Context, shopReq pkg.CreateShopReq) (string, error) {
	res, err := c.service.RegisterShop(ctx, &shop.CreateRequest{
		OwnerId:  shopReq.OwnerId,
		ShopName: shopReq.ShopName,
	})
	if err != nil {
		return "", err
	}
	return res.Shop.ShopId, err
}

func (c *ShopClient) GetShopByOwnerId(ctx context.Context, id string) ([]string, error) {
	res, err := c.service.GetShopByOwnerId(ctx, &shop.GetByOwnerIdRequest{
		OwnerId: id,
	})
	if err != nil {
		return nil, err
	}
	shopIds := []string{}
	for _, shop := range res.Shops {
		shopIds = append(shopIds, shop.ShopId)
	}
	return shopIds, nil
}
