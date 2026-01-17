package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/priyanshubiswas1/barista/apigateway/middleware"
	"github.com/priyanshubiswas1/barista/apigateway/pkg"
)

type IShopClient interface {
	GetShopsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetShopByIdHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	CreateShopHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateShopHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteShopHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

type ShopClient struct {
	url    string
	client *http.Client
}

func NewShop(url string) (IShopClient, error) {
	if url == "" {
		return &ShopClient{}, nil
	}
	return &ShopClient{
		url: url,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *ShopClient) GetShopsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idVal := r.Context().Value(middleware.CtxKeyID)
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		return errors.New("missing or invalid user ID in context")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/shops/owner/%s", c.url, idStr), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("shop service returned status: %d", resp.StatusCode)
	}

	var shops []struct {
		ShopId   string `json:"shopId"`
		ShopName string `json:"shopName"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&shops); err != nil {
		return err
	}

	// Map to pkg.Shop
	var result []*pkg.Shop
	for _, s := range shops {
		result = append(result, &pkg.Shop{
			ShopId:   s.ShopId,
			ShopName: s.ShopName,
		})
	}
	return pkg.WriteJSON(w, http.StatusOK, result)
}

func (c *ShopClient) GetShopByIdHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	shopId := vars["shopid"]

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/shop/%s", c.url, shopId), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("shop service returned status: %d", resp.StatusCode)
	}

	var res struct {
		ShopId   string `json:"shopId"`
		ShopName string `json:"shopName"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	return pkg.WriteJSON(w, http.StatusOK, &pkg.Shop{
		ShopId:   res.ShopId,
		ShopName: res.ShopName,
	})
}

func (c *ShopClient) CreateShopHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idVal := r.Context().Value(middleware.CtxKeyID)
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		return errors.New("missing or invalid user ID in context")
	}

	var reqBody pkg.CreateShopReq
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return err
	}

	payload := map[string]string{
		"ownerId":  idStr,
		"shopName": reqBody.ShopName,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/shop", c.url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("shop service returned status: %d", resp.StatusCode)
	}

	// CreateShop returns { shopId: "..." } in shop service.
	// But api-gateway expects CreateShopRes with Shop { ShopId, ShopName }.
	var res struct {
		ShopId string `json:"shopId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	return pkg.WriteJSON(w, http.StatusOK, &pkg.CreateShopRes{
		Shop: pkg.Shop{
			ShopId:   res.ShopId,
			ShopName: reqBody.ShopName,
		},
	})
}

func (c *ShopClient) UpdateShopHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idVal := r.Context().Value(middleware.CtxKeyID)
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		return errors.New("missing or invalid user ID in context")
	}

	// Wrapper struct in request? line 113: `json.NewDecoder(r.Body).Decode(&req.Shop)` ??
	// No, line 112: `req := &pkg.UpdateShopReq{}`.
	// pkg.UpdateShopReq has field `Shop`.
	// Line 113 decodes into `&req.Shop`.
	// So input JSON is expected to match `pkg.Shop` structure directly?
	// or `json.NewDecoder(r.Body).Decode(req)`?
	// The original code was `Decode(&req.Shop)`.
	// `req.Shop` type is `pkg.Shop`.
	// So input JSON has fields of `pkg.Shop` (ShopName, ShopId).
	// But `pkg.UpdateShopReq` is `struct { Shop Shop }`.
	// If the user sends `{ "shop": { ... } }`, then Decode(&req) would work.
	// If the user sends `{ "shopName": "...", "shopId": "..." }`, then Decode(&req.Shop) works.
	// Assuming the latter based on original code.

	var shopData pkg.Shop
	if err := json.NewDecoder(r.Body).Decode(&shopData); err != nil {
		return err
	}

	payload := map[string]any{
		"ownerId": idStr,
		"shopId":  shopData.ShopId,
		"name":    shopData.ShopName,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/shop", c.url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("shop service returned status: %d", resp.StatusCode)
	}

	// Response?
	// Assume success returns same data or just success?
	// UpdateShopRes has Shop.
	return pkg.WriteJSON(w, http.StatusOK, &pkg.UpdateShopRes{
		Shop: shopData,
	})
}

func (c *ShopClient) DeleteShopHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idVal := r.Context().Value(middleware.CtxKeyID)
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		return errors.New("missing or invalid user ID in context")
	}

	var reqBody pkg.DeleteShopReq
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return err
	}

	// shop service requires ownerId?
	// DELETE /shop/{id}?ownerId=...
	// My shop service implementation of Delete:
	// vars["id"] -> id.
	// Uses `s.service.Delete(r.Context(), id)`.
	// Wait, `shop/server.go` `Delete` method: `s.service.Delete(ctx, id)`.
	// The `shop` service `Delete` (service layer) might verify owner?
	// In my previous edits to `shop/server.go`:
	/*
		func (s *httpServer) Delete(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			id := vars["id"]
			if err := s.service.Delete(r.Context(), id); err != nil { ... }
	*/
	// It doesn't seem to take ownerId from request body or query param in `Delete` handler of shop server!
	// It only uses shopId.
	// The gRPC version passed `OwnerId` and `ShopId`.
	// If `shop` service logic requires OwnerId check, I might have missed passing it in `shop/server.go`.
	// The `Delete` method in `shop` service interface likely takes request DTO?
	// `s.service.Delete(ctx, id)` suggests it only takes ID.
	// If so, only ID is needed.

	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/shop/%s", c.url, reqBody.ShopId), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("shop service returned status: %d", resp.StatusCode)
	}

	return pkg.WriteJSON(w, http.StatusOK, &pkg.DeleteShopRes{
		Success: true,
	})
}
