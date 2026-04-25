package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	GetItemsAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateItemAvailable(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetItemScan3D(ctx context.Context, w http.ResponseWriter, r *http.Request) error

	//MENU-ITEMS
	CreateMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetMenuItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetMenuItemsAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteMenuItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteMenuItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error

	//MENU-ITEMS
	CreateCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetCategories(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetCategoriesWithItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error

	//IMAGE UPLOAD
	UploadImage(ctx context.Context, w http.ResponseWriter, r *http.Request) error
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
	contentType := r.Header.Get("Content-Type")
	
	var itemName, shopId, description, categoryId string
	var veg bool
	var price float32
	var available bool
	var imageUrl string

	// Check if request is multipart/form-data (for file uploads)
	if contentType != "" && len(contentType) >= 19 && contentType[:19] == "multipart/form-data" {
		// Parse multipart form with 10MB max memory
		err := r.ParseMultipartForm(10 << 20) // 10MB
		if err != nil {
			return err
		}

		// Extract form fields
		itemName = r.FormValue("name")
		shopId = r.FormValue("shopId")
		description = r.FormValue("description")
		categoryId = r.FormValue("categoryId")
		
		// Parse boolean values
		vegStr := r.FormValue("veg")
		veg = vegStr == "true"
		
		availableStr := r.FormValue("available")
		available = availableStr == "true" || availableStr == ""

		// Parse price
		priceStr := r.FormValue("price")
		if priceStr != "" {
			parsedPrice, err := strconv.ParseFloat(priceStr, 32)
			if err != nil {
				return fmt.Errorf("invalid price: %w", err)
			}
			price = float32(parsedPrice)
		}

		// Handle image upload if present
		file, header, err := r.FormFile("image")
		if err == nil {
			// Image file provided, upload it
			defer file.Close()

			// Read file data into bytes
			imageData, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			// Get filename
			filename := header.Filename
			if filename == "" {
				filename = "image.jpg"
			}

			// Upload image to Cloudinary via menu service
			uploadRes, err := c.service.UploadImage(ctx, &menu.UploadImageRequest{
				ImageData: imageData,
				Filename:  filename,
			})
			if err != nil {
				return fmt.Errorf("failed to upload image: %w", err)
			}
			imageUrl = uploadRes.ImageUrl
		}
	} else {
		req := pkg.CreateItemRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return err
		}
		itemName = req.ItemName
		shopId = req.ShopId
		description = req.Description
		categoryId = req.CategoryId
		veg = req.Veg
		price = req.Price
		imageUrl = req.Image
		available = req.Available
	}

	res, err := c.service.CreateItem(ctx, &menu.CreateItemRequest{
		Name:        itemName,
		ShopId:      shopId,
		Description: description,
		CategoryId:  categoryId,
		Veg:         veg,
		Price:       price,
		Image:       imageUrl,
		IsAvailable: available,
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
			Veg:         res.Item.Veg,
			Price:       res.Item.Price,
			Image:       res.Item.Image,
			Scan3DUrl:   res.Item.Scan3DUrl,
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
			Veg:         res.Item.Veg,
			Price:       res.Item.Price,
			Image:       res.Item.Image,
			Scan3DUrl:   res.Item.Scan3DUrl,
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
			Veg:         it.Veg,
			Price:       it.Price,
			Image:       it.Image,
			Scan3DUrl:   it.Scan3DUrl,
			Available:   it.IsAvailable,
		}
		items = append(items, item)
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetItemsResponse{
		Items: items,
	})
}

func (c *MenuClient) GetItemsAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	shopId := vars["shopid"]

	res, err := c.service.GetItemsAdmin(ctx, &menu.GetItemsRequest{
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
			Veg:         it.Veg,
			Price:       it.Price,
			Image:       it.Image,
			Scan3DUrl:   it.Scan3DUrl,
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
		Veg:         item.Veg,
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

func (c *MenuClient) UpdateItemAvailable(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	item := pkg.UpdateItemAvailableRequest{}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		return err
	}

	res, err := c.service.UpdateItemAvailable(ctx, &menu.UpdateItemAvailableRequest{
		ItemId:      item.ItemId,
		ShopId:      item.ShopId,
		IsAvailable: item.Available,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.UpdateItemAvailableResponse{
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
			Veg:         i.Veg,
			Price:       i.Price,
			Image:       i.Image,
			Scan3DUrl:   i.Scan3DUrl,
			Available:   i.IsAvailable,
		}
		items = append(items, item)
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetMenuItemsResponse{
		Items: items,
	})
}

func (c *MenuClient) GetMenuItemsAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	menuId := vars["menuid"]

	res, err := c.service.GetMenuItemsAdmin(ctx, &menu.GetMenuItemsRequest{
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
			Veg:         i.Veg,
			Price:       i.Price,
			Image:       i.Image,
			Scan3DUrl:   i.Scan3DUrl,
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

func (c *MenuClient) UpdateCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.Category{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.UpdateCategory(ctx, &menu.UpdateCategoryRequest{
		CategoryId: req.CategoryId,
		ShopId:     req.ShopId,
		Name:       req.Name,
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

func (c *MenuClient) GetCategoriesWithItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	shopId := vars["shopid"]

	res, err := c.service.GetCategoryWithItems(ctx, &menu.GetCategoryRequest{
		ShopId: shopId,
	})
	if err != nil {
		return err
	}
	categories := []pkg.CategoryWithItems{}
	for _, c := range res.Categories {
		items := []pkg.Item{}
		for _, i := range c.Items {
			items = append(items, pkg.Item{
				ItemId:      i.ItemId,
				ShopId:      i.ShopId,
				ItemName:    i.Name,
				Description: i.Description,
				Category:    i.Category,
				Veg:         i.Veg,
				Price:       i.Price,
				Image:       i.Image,
				Scan3DUrl:   i.Scan3DUrl,
				Available:   i.IsAvailable,
			})
		}
		category := pkg.CategoryWithItems{
			CategoryId: c.Id,
			ShopId:     c.ShopId,
			Name:       c.Name,
			Items:      items,
		}
		categories = append(categories, category)
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetCategoriesWithItemsResponse{
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

func (c *MenuClient) UploadImage(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Parse multipart form with 10MB max memory
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		return err
	}

	// Get the file from the form
	file, header, err := r.FormFile("image")
	if err != nil {
		return err
	}
	defer file.Close()

	// Read file data into bytes
	imageData, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Get filename
	filename := header.Filename
	if filename == "" {
		filename = "image.jpg"
	}

	// Call menu service
	res, err := c.service.UploadImage(ctx, &menu.UploadImageRequest{
		ImageData: imageData,
		Filename:  filename,
	})
	if err != nil {
		return err
	}

	return pkg.WriteJSON(w, http.StatusOK, pkg.UploadImageResponse{
		ImageUrl: res.ImageUrl,
	})
}

func (c *MenuClient) GetItemScan3D(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	itemId := vars["itemid"]

	res, err := c.service.GetItem(ctx, &menu.GetItemRequest{ItemId: itemId})
	if err != nil {
		return err
	}
	scanURL := strings.TrimSpace(res.GetItem().GetScan3DUrl())
	if scanURL == "" {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("3d scan not available"))
		return nil
	}
	if !strings.HasPrefix(scanURL, "s3://") {
		return fmt.Errorf("unsupported scan3d url: %s", scanURL)
	}
	trimmed := strings.TrimPrefix(scanURL, "s3://")
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid s3 url: %s", scanURL)
	}
	bucket := parts[0]
	key := parts[1]

	endpoint := os.Getenv("GARAGE_ENDPOINT")
	region := os.Getenv("GARAGE_REGION")
	accessKey := os.Getenv("GARAGE_ACCESS_KEY")
	secretKey := os.Getenv("GARAGE_SECRET_KEY")
	if endpoint == "" || region == "" || accessKey == "" || secretKey == "" {
		return fmt.Errorf("garage env not configured in api-gateway")
	}

	awsCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
				if service == s3.ServiceID {
					return aws.Endpoint{URL: endpoint, HostnameImmutable: true}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			},
		)),
	)
	if err != nil {
		return err
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	obj, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return err
	}
	defer obj.Body.Close()

	w.Header().Set("Content-Type", "model/gltf-binary")
	w.Header().Set("Content-Disposition", "inline; filename=\"model.glb\"")
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, obj.Body)
	return err
}
