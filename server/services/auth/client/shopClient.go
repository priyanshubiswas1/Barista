package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/priyanshubiswas1/barista/authservice/pkg"
)

type IShopClient interface {
	RegisterShop(ctx context.Context, shop pkg.CreateShopReq) (string, error)
	GetShopByOwnerId(ctx context.Context, id string) ([]string, error)
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

func (c *ShopClient) RegisterShop(ctx context.Context, shopReq pkg.CreateShopReq) (string, error) {
	if c.url == "" {
		return "", errors.New("shop service URL not configured")
	}

	payload := map[string]string{
		"ownerId":  shopReq.OwnerId,
		"shopName": shopReq.ShopName,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/shop", c.url), bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("shop service returned status: %d", resp.StatusCode)
	}

	var res struct {
		ShopId string `json:"shopId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.ShopId, nil
}

func (c *ShopClient) GetShopByOwnerId(ctx context.Context, id string) ([]string, error) {
	if c.url == "" {
		return nil, errors.New("shop service URL not configured")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/shops/owner/%s", c.url, id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("shop service returned status: %d", resp.StatusCode)
	}

	var res []struct {
		ShopId string `json:"shopId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	shopIds := []string{}
	for _, s := range res {
		shopIds = append(shopIds, s.ShopId)
	}
	return shopIds, nil
}
