package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/priyanshubiswas1/barista/apigateway/middleware"
	"github.com/priyanshubiswas1/barista/apigateway/pkg"
	shop "github.com/priyanshubiswas1/baristacontracts/gen/shop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IShopClient interface {
	GetShopsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetShopByIdHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	CreateShopHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateShopHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteShopHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
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

func (c *ShopClient) GetShopsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idVal := r.Context().Value(middleware.CtxKeyID)
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		return errors.New("missing or invalid user ID in context")
	}
	res, err := c.service.GetShopByOwnerId(ctx, &shop.GetByOwnerIdRequest{
		OwnerId: idStr,
	})
	if err != nil {
		return err
	}
	shops := []*pkg.Shop{}
	for _, s := range res.Shops {
		shop := &pkg.Shop{
			ShopId:   s.ShopId,
			ShopName: s.ShopName,
		}
		shops = append(shops, shop)
	}
	return pkg.WriteJSON(w, http.StatusOK, shops)
}

func (c *ShopClient) GetShopByIdHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	shopId := vars["shopid"]
	res, err := c.service.GetShopById(ctx, &shop.GetByShopIdRequest{
		ShopId: shopId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, &pkg.Shop{
		ShopId:   res.Shop.ShopId,
		ShopName: res.Shop.ShopName,
	})
}

func (c *ShopClient) CreateShopHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idVal := r.Context().Value(middleware.CtxKeyID)
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		return errors.New("missing or invalid user ID in context")
	}
	req := &pkg.CreateShopReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.RegisterShop(ctx, &shop.CreateRequest{
		OwnerId:  idStr,
		ShopName: req.ShopName,
	})
	if err != nil {
		return err
	}
	shop := &pkg.Shop{
		ShopName: res.Shop.ShopName,
		ShopId:   res.Shop.ShopId,
	}
	return pkg.WriteJSON(w, http.StatusOK, &pkg.CreateShopRes{
		Shop: *shop,
	})
}

func (c *ShopClient) UpdateShopHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idVal := r.Context().Value(middleware.CtxKeyID)
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		return errors.New("missing or invalid user ID in context")
	}
	req := &pkg.UpdateShopReq{}
	if err := json.NewDecoder(r.Body).Decode(&req.Shop); err != nil {
		return err
	}
	shopReqBody := &shop.Shop{
		ShopId:   req.Shop.ShopId,
		ShopName: req.Shop.ShopName,
	}
	res, err := c.service.UpdateShop(ctx, &shop.UpdateRequest{
		OwnerId: idStr,
		Shop:    shopReqBody,
	})
	if err != nil {
		return err
	}
	shopResBody := &pkg.Shop{
		ShopId:   res.Shop.ShopId,
		ShopName: res.Shop.ShopName,
	}
	return pkg.WriteJSON(w, http.StatusOK, &pkg.UpdateShopRes{
		Shop: *shopResBody,
	})
}

func (c *ShopClient) DeleteShopHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idVal := r.Context().Value(middleware.CtxKeyID)
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		return errors.New("missing or invalid user ID in context")
	}
	req := &pkg.DeleteShopReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.DeleteShop(ctx, &shop.DeleteRequest{
		OwnerId: idStr,
		ShopId:  req.ShopId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, &pkg.DeleteShopRes{
		Success: res.Success,
	})
}
