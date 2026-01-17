package menu

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/priyanshubiswas1/barista/menuservice/pkg"
)

type httpServer struct {
	service IService
}

func ListenHTTP(service IService, port string) error {
	router := mux.NewRouter()
	s := &httpServer{service: service}

	// Health checks
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	router.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	}).Methods("GET")

	// API Routes
	// Menu
	router.HandleFunc("/menu", s.CreateMenu).Methods("POST")
	router.HandleFunc("/menus/{shopid}", s.GetMenus).Methods("GET")
	router.HandleFunc("/menu/{menuid}", s.SetMenu).Methods("PUT")
	router.HandleFunc("/menu/{menuid}", s.DeleteMenu).Methods("DELETE")

	// Items
	router.HandleFunc("/item", s.CreateItem).Methods("POST")
	router.HandleFunc("/item/{itemid}", s.GetItem).Methods("GET")
	router.HandleFunc("/items/{shopid}", s.GetItems).Methods("GET")
	router.HandleFunc("/items/admin/{shopid}", s.GetItemsAdmin).Methods("GET")
	router.HandleFunc("/item", s.UpdateItem).Methods("PUT")
	router.HandleFunc("/item/available", s.UpdateItemAvailable).Methods("PUT")
	router.HandleFunc("/item/{itemid}", s.DeleteItem).Methods("DELETE")

	// Menu Items
	router.HandleFunc("/menu-item", s.CreateMenuItem).Methods("POST")
	router.HandleFunc("/menu-item/{menuid}", s.GetMenuItems).Methods("GET")
	router.HandleFunc("/menu-item/admin/{menuid}", s.GetMenuItemsAdmin).Methods("GET")
	router.HandleFunc("/menu-item", s.DeleteMenuItem).Methods("DELETE") // Params in body or query?
	router.HandleFunc("/menu-items/{menuid}", s.DeleteMenuItems).Methods("DELETE")

	// Categories
	router.HandleFunc("/category", s.CreateCategory).Methods("POST")
	router.HandleFunc("/category/{shopid}", s.GetCategories).Methods("GET")
	router.HandleFunc("/category-items/{shopid}", s.GetCategoryWithItems).Methods("GET")
	router.HandleFunc("/category", s.UpdateCategory).Methods("PUT")
	router.HandleFunc("/category/{id}", s.DeleteCategory).Methods("DELETE")

	// Image Upload
	router.HandleFunc("/upload-image", s.UploadImage).Methods("POST")

	log.Printf("menu service listening on HTTP %s", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}

func (s *httpServer) CreateMenu(w http.ResponseWriter, r *http.Request) {
	var req pkg.CreateMenuReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := s.service.CreateMenu(r.Context(), req.ShopId, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, pkg.Menu{
		MenuId: id,
		ShopId: req.ShopId,
		Name:   req.Name,
	})
}

func (s *httpServer) GetMenus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shopId := vars["shopid"]
	menus, err := s.service.GetMenus(r.Context(), shopId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"menus": menus})
}

func (s *httpServer) SetMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuId := vars["menuid"]
	success, err := s.service.SetMenu(r.Context(), menuId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": success})
}

func (s *httpServer) DeleteMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuId := vars["menuid"]
	success, err := s.service.DeleteMenu(r.Context(), menuId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": success})
}

func (s *httpServer) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req pkg.ItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := s.service.CreateItem(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"item": item})
}

func (s *httpServer) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemId := vars["itemid"]
	item, err := s.service.GetItem(r.Context(), itemId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"item": item})
}

func (s *httpServer) GetItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shopId := vars["shopid"]
	items, err := s.service.GetItems(r.Context(), shopId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (s *httpServer) GetItemsAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shopId := vars["shopid"]
	items, err := s.service.GetItemsAdmin(r.Context(), shopId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (s *httpServer) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var req pkg.ItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	success, err := s.service.UpdateItem(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": success})
}

func (s *httpServer) UpdateItemAvailable(w http.ResponseWriter, r *http.Request) {
	var req pkg.UpdateItemAvailableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	success, err := s.service.UpdateItemAvailable(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": success})
}

func (s *httpServer) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemId := vars["itemid"]
	success, err := s.service.DeleteItem(r.Context(), itemId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": success})
}

func (s *httpServer) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	var req pkg.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	menuItem, err := s.service.CreateMenuItem(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Warning: The service returns MenuItem, make sure to json properly
	writeJSON(w, http.StatusOK, map[string]any{"menuItem": menuItem})
}

func (s *httpServer) GetMenuItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuId := vars["menuid"]
	items, err := s.service.GetMenuItems(r.Context(), menuId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (s *httpServer) GetMenuItemsAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuId := vars["menuid"]
	items, err := s.service.GetMenuItemsAdmin(r.Context(), menuId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (s *httpServer) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	// Attempt to get from body first (uncommon for DELETE)
	// Or query params?
	var req struct {
		MenuId string `json:"menuId"`
		ItemId string `json:"itemId"`
	}
	// Check query params first
	req.MenuId = r.URL.Query().Get("menuId")
	req.ItemId = r.URL.Query().Get("itemId")

	if req.MenuId == "" || req.ItemId == "" {
		// Try body
		json.NewDecoder(r.Body).Decode(&req)
	}

	if req.MenuId == "" || req.ItemId == "" {
		http.Error(w, "missing menuId or itemId", http.StatusBadRequest)
		return
	}

	success, err := s.service.DeleteMenuItem(r.Context(), req.MenuId, req.ItemId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": success})
}

func (s *httpServer) DeleteMenuItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuId := vars["menuid"]
	success, err := s.service.DeleteMenuItems(r.Context(), menuId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": success})
}

func (s *httpServer) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req pkg.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	category, err := s.service.CreateCategory(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"category": category})
}

func (s *httpServer) GetCategories(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shopId := vars["shopid"]
	categories, err := s.service.GetCategories(r.Context(), shopId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"categories": categories})
}

func (s *httpServer) GetCategoryWithItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shopId := vars["shopid"]
	categories, err := s.service.GetCategoriesWithItem(r.Context(), shopId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"categories": categories})
}

func (s *httpServer) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var req pkg.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	category, err := s.service.UpdateCategory(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"category": category})
}

func (s *httpServer) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// Need shopId?
	shopId := r.URL.Query().Get("shopId")
	// If missing, try body?
	if shopId == "" {
		var req struct {
			ShopId string `json:"shopId"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		shopId = req.ShopId
	}

	success, err := s.service.DeleteCategory(r.Context(), id, shopId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": success})
}

func (s *httpServer) UploadImage(w http.ResponseWriter, r *http.Request) {
	// Parse 10MB
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "image file required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url, err := s.service.UploadImage(r.Context(), data, header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"imageUrl": url})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
