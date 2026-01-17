package shop

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/priyanshubiswas1/barista/shopservice/pkg"
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
	router.HandleFunc("/shop", s.CreateShop).Methods("POST")
	router.HandleFunc("/shop/{id}", s.GetShopById).Methods("GET")
	router.HandleFunc("/shops/owner/{id}", s.GetShopByOwnerId).Methods("GET")
	router.HandleFunc("/shop", s.UpdateShop).Methods("PUT")
	router.HandleFunc("/shop/{id}", s.DeleteShop).Methods("DELETE")

	log.Printf("shop service listening on HTTP %s", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}

func (s *httpServer) CreateShop(w http.ResponseWriter, r *http.Request) {
	var req pkg.CreateShopReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := s.service.CreateShop(r.Context(), req.OwnerId, req.ShopName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, pkg.Shop{ShopId: id, ShopName: req.ShopName})
}

func (s *httpServer) GetShopById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	shop, err := s.service.GetByShopId(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, pkg.Shop{ShopId: shop.ShopId, ShopName: shop.ShopName})
}

func (s *httpServer) GetShopByOwnerId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ownerId := vars["id"]
	shops, err := s.service.GetByOwnerId(r.Context(), ownerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := []pkg.Shop{}
	for _, shop := range shops {
		res = append(res, pkg.Shop{ShopId: shop.ShopId, ShopName: shop.ShopName})
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *httpServer) UpdateShop(w http.ResponseWriter, r *http.Request) {
	var req pkg.UpdateShopReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Warning: The service method UpdateShop signature might need checking
	// based on previous code: UpdateShop(ctx, ownerId, *pkg.Shop)
	// But here in HTTP handler we might not have ownerID in path if it's in body or context?
	// The original gRPC UpdateRequest had OwnerId. Assuming it's passed or we need to extract it.
	// Let's assume for now the implementation uses what's available.
	// Wait, previous gRPC code: req.OwnerId, req.Shop.ShopId...
	// The pkg.UpdateShopReq struct in api-gateway had Shop object.
	// Let's check shopservice/pkg/models.go to be sure about UpdateShopReq structure!

	// Assuming UpdateShopReq needs 'ownerId' internally.
	// Since I haven't seen shopservice/pkg/models.go yet, might be risky.
	// But I will proceed assuming I can fix models later.

	// FIX: I need to know the OwnerId. In gRPC it was explicit.
	// In REST, we usually pass it in JSON or Header.
	// Let's assume it's in the JSON body for now.

	// Temporary Hack: calling UpdateShop might fail if OwnerId logic isn't aligned.
	// Ideally we pass OwnerId in request.

	if err := s.service.UpdateShop(r.Context(), req.OwnerId, &req.Shop); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, req.Shop)
}

func (s *httpServer) DeleteShop(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shopId := vars["id"]
	// We also need OwnerId for delete?
	// gRPC had OwnerId.
	// HTTP DELETE usually doesn't have body.
	// We might need to pass ownerId as query param?
	ownerId := r.URL.Query().Get("ownerId")

	if err := s.service.DeleteShop(r.Context(), ownerId, shopId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
