package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/priyanshubiswas1/barista/apigateway/handler"
	"github.com/priyanshubiswas1/barista/apigateway/middleware"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

type APIServer struct {
	authClientUrl string
	shopClientUrl string
	menuClientUrl string
}

func NewAPIServer(authClient, shopClientUrl, menuClientUrl string) (*APIServer, error) {
	return &APIServer{
		authClientUrl: authClient,
		shopClientUrl: shopClientUrl,
		menuClientUrl: menuClientUrl,
	}, nil
}

func (s *APIServer) ListenServer(port string) error {
	authClient, err := handler.NewClient(s.authClientUrl)
	if err != nil {
		return err
	}
	shopClient, err := handler.NewShop(s.shopClientUrl)
	if err != nil {
		return err
	}
	menuClient, err := handler.NewMenu(s.menuClientUrl)
	if err != nil {
		return err
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/auth/register", makeHttpHandleFunc(authClient.RegisterHandler)).Methods("POST")
	router.HandleFunc("/api/auth/login", makeHttpHandleFunc(authClient.LoginHandler)).Methods("POST")

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	//USER
	protected.HandleFunc("/api/auth/user", makeHttpHandleFunc(authClient.GetUserHandler)).Methods("GET")
	protected.HandleFunc("/api/auth/user", makeHttpHandleFunc(authClient.UpdateUserHandler)).Methods("PUT")
	protected.HandleFunc("/api/auth/user", makeHttpHandleFunc(authClient.DeleteUserHandler)).Methods("DELETE")

	//SHOP
	protected.HandleFunc("/api/shop/{shopid}", makeHttpHandleFunc(shopClient.GetShopByIdHandler)).Methods("GET")
	protected.HandleFunc("/api/shops", makeHttpHandleFunc(shopClient.GetShopsHandler)).Methods("GET")
	protected.HandleFunc("/api/shop", makeHttpHandleFunc(shopClient.CreateShopHandler)).Methods("POST")
	protected.HandleFunc("/api/shop", makeHttpHandleFunc(shopClient.UpdateShopHandler)).Methods("PUT")
	protected.HandleFunc("/api/shop", makeHttpHandleFunc(shopClient.DeleteShopHandler)).Methods("DELETE")

	//MENU
	protected.HandleFunc("/api/menu", makeHttpHandleFunc(menuClient.CreateMenu)).Methods("POST")
	protected.HandleFunc("/api/menus/{shopid}", makeHttpHandleFunc(menuClient.GetMenus)).Methods("GET")
	protected.HandleFunc("/api/menu/{menuid}", makeHttpHandleFunc(menuClient.SetMenu)).Methods("PUT")
	protected.HandleFunc("/api/menu/{menuid}", makeHttpHandleFunc(menuClient.DeleteMenu)).Methods("DELETE")

	//ITEMS
	protected.HandleFunc("/api/item", makeHttpHandleFunc(menuClient.CreateItem)).Methods("POST")
	protected.HandleFunc("/api/item/{itemid}", makeHttpHandleFunc(menuClient.GetItem)).Methods("GET")
	protected.HandleFunc("/api/items/{shopid}", makeHttpHandleFunc(menuClient.GetItems)).Methods("GET")
	protected.HandleFunc("/api/item/{itemid}", makeHttpHandleFunc(menuClient.UpdateItem)).Methods("PUT")
	protected.HandleFunc("/api/item/{itemid}", makeHttpHandleFunc(menuClient.DeleteItem)).Methods("DELETE")

	//MENU-ITEMS
	protected.HandleFunc("/api/menu-item", makeHttpHandleFunc(menuClient.CreateMenuItem)).Methods("POST")
	protected.HandleFunc("/api/menu-item/{menuid}", makeHttpHandleFunc(menuClient.GetMenuItems)).Methods("GET")
	protected.HandleFunc("/api/menu-item", makeHttpHandleFunc(menuClient.DeleteMenuItem)).Methods("DELETE")
	protected.HandleFunc("/api/menu-item/{menuid}", makeHttpHandleFunc(menuClient.DeleteMenuItems)).Methods("DELETE")

	return http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}

func makeHttpHandleFunc(apiFunc APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		if err := apiFunc(ctx, w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		}
	}
}

func writeJSON(w http.ResponseWriter, s int, data any) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(data)
}
