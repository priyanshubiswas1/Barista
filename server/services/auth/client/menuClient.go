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

type IMenuClient interface {
	CreateMenu(ctx context.Context, menuReq pkg.CreateMenuReq) (string, error)
	GetMenus(ctx context.Context, shopId string) ([]string, error)
}

type MenuClient struct {
	url    string
	client *http.Client
}

func NewMenu(url string) (IMenuClient, error) {
	if url == "" {
		return &MenuClient{}, nil
	}
	return &MenuClient{
		url: url,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *MenuClient) CreateMenu(ctx context.Context, menuReq pkg.CreateMenuReq) (string, error) {
	if c.url == "" {
		return "", errors.New("menu service URL not configured")
	}

	payload := map[string]string{
		"shopId": menuReq.ShopId,
		"name":   menuReq.MenuName,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/menu", c.url), bytes.NewBuffer(body))
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
		return "", fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res struct {
		MenuId string `json:"menuId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.MenuId, nil
}

func (c *MenuClient) GetMenus(ctx context.Context, shopId string) ([]string, error) {
	if c.url == "" {
		return nil, errors.New("menu service URL not configured")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/menus/%s", c.url, shopId), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res struct {
		Menus []struct {
			MenuId string `json:"menuId"`
		} `json:"menus"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	menuIds := []string{}
	for _, m := range res.Menus {
		menuIds = append(menuIds, m.MenuId)
	}
	return menuIds, nil
}
