package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/priyanshubiswas1/barista/apigateway/pkg"
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
	GetItemsAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateItemAvailable(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error

	//MENU-ITEMS
	CreateMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetMenuItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetMenuItemsAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteMenuItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error

	//CATEGORIES
	CreateCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetCategories(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetCategoriesWithItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error

	//IMAGE UPLOAD
	UploadImage(ctx context.Context, w http.ResponseWriter, r *http.Request) error
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
			// Increased to 60s to handle cold starts (and large image uploads)
			Timeout: 60 * time.Second,
		},
	}, nil
}

func (c *MenuClient) CreateMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req pkg.CreateMenuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/menu", c.url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res pkg.Menu
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	return pkg.WriteJSON(w, http.StatusOK, pkg.CreateMenuResponse{
		Menu: res,
	})
}

func (c *MenuClient) GetMenus(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	shopId := vars["shopid"]

	httpReq, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/menus/%s", c.url, shopId), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res struct {
		Menus []pkg.Menu `json:"menus"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	return pkg.WriteJSON(w, http.StatusOK, pkg.GetMenusResponse{
		Menus: res.Menus,
	})
}

func (c *MenuClient) SetMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuId := vars["menuid"]

	httpReq, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/menu/%s", c.url, menuId), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string]bool
	json.NewDecoder(resp.Body).Decode(&res)

	return pkg.WriteJSON(w, http.StatusOK, pkg.SetMenuResponse{
		Success: res["success"],
	})
}

func (c *MenuClient) DeleteMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuId := vars["menuid"]

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/menu/%s", c.url, menuId), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string]bool
	json.NewDecoder(resp.Body).Decode(&res)

	return pkg.WriteJSON(w, http.StatusOK, pkg.DeleteMenuResponse{
		Success: res["success"],
	})
}

// CreateItem - Supports Multipart
func (c *MenuClient) CreateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	contentType := r.Header.Get("Content-Type")

	// If it's multipart, we need to forward it as multipart to menu service?
	// OR construct a new multipart request.
	if contentType != "" && len(contentType) >= 19 && contentType[:19] == "multipart/form-data" {
		// New Multipart Request Buffer
		var b bytes.Buffer
		wMw := multipart.NewWriter(&b)

		// Parse Incoming Form
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			return err
		}

		// Copy fields
		for key, val := range r.MultipartForm.Value {
			// val is []string, usually take first
			if len(val) > 0 {
				wMw.WriteField(key, val[0])
			}
		}

		// Copy Files (Image)
		for key, files := range r.MultipartForm.File {
			for _, fileHeader := range files {
				file, err := fileHeader.Open()
				if err != nil {
					return err
				}
				fw, err := wMw.CreateFormFile(key, fileHeader.Filename)
				if err != nil {
					file.Close()
					return err
				}
				io.Copy(fw, file)
				file.Close()
			}
		}
		wMw.Close()

		httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/item", c.url), &b)
		if err != nil {
			return err
		}
		httpReq.Header.Set("Content-Type", wMw.FormDataContentType())

		resp, err := c.client.Do(httpReq)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
		}

		var res map[string]pkg.Item // { "item": ... }
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return err
		}
		return pkg.WriteJSON(w, http.StatusOK, pkg.CreateItemResponse{
			Item: res["item"],
		})

	} else {
		// JSON
		var req pkg.CreateItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return err
		}

		body, _ := json.Marshal(req)
		httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/item", c.url), bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := c.client.Do(httpReq)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
		}

		var res map[string]pkg.Item
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return err
		}
		return pkg.WriteJSON(w, http.StatusOK, pkg.CreateItemResponse{
			Item: res["item"],
		})
	}
}

func (c *MenuClient) GetItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	itemId := vars["itemid"]

	httpReq, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/item/%s", c.url, itemId), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string]pkg.Item
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetItemResponse{
		Item: res["item"],
	})
}

func (c *MenuClient) GetItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	shopId := vars["shopid"]

	httpReq, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/items/%s", c.url, shopId), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string][]pkg.Item
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetItemsResponse{
		Items: res["items"],
	})
}

func (c *MenuClient) GetItemsAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	shopId := vars["shopid"]

	httpReq, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/items/admin/%s", c.url, shopId), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string][]pkg.Item
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetItemsResponse{
		Items: res["items"],
	})
}

func (c *MenuClient) UpdateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req pkg.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/item", c.url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string]bool
	json.NewDecoder(resp.Body).Decode(&res)

	return pkg.WriteJSON(w, http.StatusOK, pkg.UpdateItemResponse{
		Success: res["success"],
	})
}

func (c *MenuClient) UpdateItemAvailable(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req pkg.UpdateItemAvailableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/item/available", c.url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string]bool
	json.NewDecoder(resp.Body).Decode(&res)

	return pkg.WriteJSON(w, http.StatusOK, pkg.UpdateItemAvailableResponse{
		Success: res["success"],
	})
}

func (c *MenuClient) DeleteItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	itemId := vars["itemid"]

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/item/%s", c.url, itemId), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string]bool
	json.NewDecoder(resp.Body).Decode(&res)

	return pkg.WriteJSON(w, http.StatusOK, pkg.DeleteItemResponse{
		Success: res["success"],
	})
}

func (c *MenuClient) CreateMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req pkg.CreateMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/menu-item", c.url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string]pkg.MenuItem
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	return pkg.WriteJSON(w, http.StatusOK, pkg.CreateMenuItemResponse{
		MenuItem: res["menuItem"],
	})
}

func (c *MenuClient) GetMenuItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuId := vars["menuid"]

	httpReq, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/menu-item/%s", c.url, menuId), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string][]pkg.Item
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetMenuItemsResponse{
		Items: res["items"],
	})
}

func (c *MenuClient) GetMenuItemsAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuId := vars["menuid"]

	httpReq, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/menu-item/admin/%s", c.url, menuId), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string][]pkg.Item
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetMenuItemsResponse{
		Items: res["items"],
	})
}

func (c *MenuClient) DeleteMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req pkg.DeleteMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Uses Query Params in `menu/server.go` implementation?
	// Ah, my logic in `DeleteMenuItem` (menu service) logic checked Query params OR Body.
	// So I can send Body or Query params.
	// I'll send Query params as a safer bet (or body).
	// Let's send Query params.

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/menu-item", c.url), nil)
	if err != nil {
		return err
	}
	q := httpReq.URL.Query()
	q.Add("menuId", req.MenuId)
	q.Add("itemId", req.ItemId)
	httpReq.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string]bool
	json.NewDecoder(resp.Body).Decode(&res)

	return pkg.WriteJSON(w, http.StatusOK, pkg.DeleteMenuItemResponse{
		Success: res["success"],
	})
}

func (c *MenuClient) DeleteMenuItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuId := vars["menuid"]

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/menu-items/%s", c.url, menuId), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string]bool
	json.NewDecoder(resp.Body).Decode(&res)

	return pkg.WriteJSON(w, http.StatusOK, pkg.DeleteMenuItemsResponse{
		Success: res["success"],
	})
}

func (c *MenuClient) CreateCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req pkg.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/category", c.url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string]pkg.Category
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.CreateCategoryResponse{
		Category: res["category"],
	})
}

func (c *MenuClient) UpdateCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req pkg.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/category", c.url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string]pkg.Category
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.CreateCategoryResponse{
		Category: res["category"],
	})
}

func (c *MenuClient) GetCategories(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	shopId := vars["shopid"]

	httpReq, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/category/%s", c.url, shopId), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string][]pkg.Category
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetCategoriesResponse{
		Categories: res["categories"],
	})
}

func (c *MenuClient) GetCategoriesWithItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	shopId := vars["shopid"]

	httpReq, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/category-items/%s", c.url, shopId), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string][]pkg.CategoryWithItems
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetCategoriesWithItemsResponse{
		Categories: res["categories"],
	})
}

func (c *MenuClient) DeleteCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	// Needs shopId
	shopId := r.URL.Query().Get("shopId")
	if shopId == "" {
		// Read body
		var req struct {
			ShopId string `json:"shopId"`
		}
		// Careful consuming body if query not present
		json.NewDecoder(r.Body).Decode(&req)
		shopId = req.ShopId
	}

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/category/%s", c.url, id), nil)
	if err != nil {
		return err
	}
	// Add query param for shopId
	q := httpReq.URL.Query()
	q.Add("shopId", shopId)
	httpReq.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("menu service returned status: %d", resp.StatusCode)
	}

	var res map[string]bool
	json.NewDecoder(resp.Body).Decode(&res)

	return pkg.WriteJSON(w, http.StatusOK, pkg.DeleteCategoryResponse{
		Success: res["success"],
	})
}

func (c *MenuClient) UploadImage(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Need to forward multipart request
	// r.Body is the multipart stream.
	// r.Header["Content-Type"] gives boundary.

	// We can't just forward r.Body because we need to construct a new request to target URL.
	// But we can pass r.Body as the body of new request, and copy Content-Type header.

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/upload-image", c.url), r.Body)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", r.Header.Get("Content-Type"))
	httpReq.ContentLength = r.ContentLength // Important for multipart

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// read error from body
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("menu service returned status: %d, body: %s", resp.StatusCode, string(b))
	}

	var res map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	return pkg.WriteJSON(w, http.StatusOK, pkg.UploadImageResponse{
		ImageUrl: res["imageUrl"],
	})
}
