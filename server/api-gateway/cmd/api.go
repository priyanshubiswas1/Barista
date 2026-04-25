package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/priyanshubiswas1/barista/apigateway/handler"
	"github.com/priyanshubiswas1/barista/apigateway/middleware"
	"github.com/rs/cors"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

type APIServer struct {
	authClientUrl    string
	shopClientUrl    string
	menuClientUrl    string
	tableClientUrl   string
	sessionClientUrl string
	orderClientUrl   string
	paymentClientUrl string
	historyClientUrl string
	allowedOrigins   string
}

func NewAPIServer(
	authClient,
	shopClientUrl,
	menuClientUrl,
	tableClientUrl,
	sessionClientUrl,
	orderClientUrl,
	paymentClientUrl,
	historyClientUrl,
	allowedOrigins string,
) (*APIServer, error) {
	return &APIServer{
		authClientUrl:    authClient,
		shopClientUrl:    shopClientUrl,
		menuClientUrl:    menuClientUrl,
		tableClientUrl:   tableClientUrl,
		sessionClientUrl: sessionClientUrl,
		orderClientUrl:   orderClientUrl,
		paymentClientUrl: paymentClientUrl,
		historyClientUrl: historyClientUrl,
		allowedOrigins:   allowedOrigins,
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
	tableClient, err := handler.NewTable(s.tableClientUrl)
	if err != nil {
		return err
	}
	sessionClient, err := handler.NewSession(s.sessionClientUrl)
	if err != nil {
		return err
	}
	orderClient, err := handler.NewOrder(s.orderClientUrl)
	if err != nil {
		return err
	}
	paymentClient, err := handler.NewPayment(s.paymentClientUrl)
	if err != nil {
		return err
	}
	historyClient, err := handler.NewHistory(s.historyClientUrl)
	if err != nil {
		return err
	}

	router := mux.NewRouter()

	// Parse allowed origins from comma-separated string
	origins := strings.Split(s.allowedOrigins, ",")
	allowedOriginsList := make([]string, 0, len(origins))
	for _, origin := range origins {
		trimmedOrigin := strings.TrimSpace(origin)
		if trimmedOrigin != "" {
			allowedOriginsList = append(allowedOriginsList, trimmedOrigin)
		}
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOriginsList,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Session"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	router.HandleFunc("/api/auth/onboard", makeHttpHandleFunc(authClient.RegisterWithShopMenuHandler)).Methods("POST")
	router.HandleFunc("/api/auth/register", makeHttpHandleFunc(authClient.RegisterHandler)).Methods("POST")
	router.HandleFunc("/api/auth/login", makeHttpHandleFunc(authClient.LoginHandler)).Methods("POST")

	// unprotected SESSION paths
	router.HandleFunc("/api/items/{shopid}", makeHttpHandleFunc(menuClient.GetItems)).Methods("GET")
	router.HandleFunc("/api/item/{itemid}/scan3d", makeHttpHandleFunc(menuClient.GetItemScan3D)).Methods("GET")
	router.HandleFunc("/api/sessions", makeHttpHandleFunc(sessionClient.GetActiveSession)).Methods("GET").Queries("shopid", "{shopid}", "tableid", "{tableid}")
	router.HandleFunc("/api/sessions/send-otp", makeHttpHandleFunc(sessionClient.SendSessionOTP)).Methods("POST")
	router.HandleFunc("/api/sessions", makeHttpHandleFunc(sessionClient.CreateSession)).Methods("POST")
	router.HandleFunc("/api/sessions/join", makeHttpHandleFunc(sessionClient.JoinSession)).Methods("POST")
	router.HandleFunc("/api/customer/table/{tableid}", makeHttpHandleFunc(tableClient.GetTable)).Methods("GET")
	router.HandleFunc("/api/order-customer", makeHttpHandleFunc(orderClient.CreateOrderCustomer)).Methods("POST")

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	// protected.Use(middleware.SessionMiddleware)

	session := router.PathPrefix("/").Subrouter()
	session.Use(middleware.SessionMiddleware)

	//USER
	protected.HandleFunc("/api/auth/verify", http.HandlerFunc(VerifyHandler)).Methods("GET")
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
	protected.HandleFunc("/api/items/{shopid}/admin", makeHttpHandleFunc(menuClient.GetItemsAdmin)).Methods("GET")
	protected.HandleFunc("/api/item", makeHttpHandleFunc(menuClient.UpdateItem)).Methods("PUT")
	protected.HandleFunc("/api/item-available", makeHttpHandleFunc(menuClient.UpdateItemAvailable)).Methods("PUT")
	protected.HandleFunc("/api/item/{itemid}", makeHttpHandleFunc(menuClient.DeleteItem)).Methods("DELETE")

	//MENU-ITEMS
	protected.HandleFunc("/api/menu-item", makeHttpHandleFunc(menuClient.CreateMenuItem)).Methods("POST")
	protected.HandleFunc("/api/menu-item/{menuid}", makeHttpHandleFunc(menuClient.GetMenuItems)).Methods("GET")
	protected.HandleFunc("/api/menu-item", makeHttpHandleFunc(menuClient.DeleteMenuItem)).Methods("DELETE")
	protected.HandleFunc("/api/menu-item/{menuid}", makeHttpHandleFunc(menuClient.DeleteMenuItems)).Methods("DELETE")

	//CATEGORIES
	protected.HandleFunc("/api/category", makeHttpHandleFunc(menuClient.CreateCategory)).Methods("POST")
	protected.HandleFunc("/api/category", makeHttpHandleFunc(menuClient.UpdateCategory)).Methods("PUT")
	protected.HandleFunc("/api/category/{shopid}", makeHttpHandleFunc(menuClient.GetCategories)).Methods("GET")
	protected.HandleFunc("/api/category-items/{shopid}", makeHttpHandleFunc(menuClient.GetCategoriesWithItems)).Methods("GET")
	protected.HandleFunc("/api/category", makeHttpHandleFunc(menuClient.DeleteCategory)).Methods("DELETE")

	//IMAGE UPLOAD
	protected.HandleFunc("/api/upload-image", makeHttpHandleFunc(menuClient.UploadImage)).Methods("POST")

	//TABLES
	protected.HandleFunc("/api/table", makeHttpHandleFunc(tableClient.CreateTable)).Methods("POST")
	protected.HandleFunc("/api/tables/{shopid}", makeHttpHandleFunc(tableClient.GetTables)).Methods("GET")
	protected.HandleFunc("/api/table/{tableid}", makeHttpHandleFunc(tableClient.GetTable)).Methods("GET")
	protected.HandleFunc("/api/table", makeHttpHandleFunc(tableClient.UpdateTable)).Methods("PUT")
	protected.HandleFunc("/api/table", makeHttpHandleFunc(tableClient.DeleteTable)).Methods("DELETE")

	//SESSIONS
	protected.HandleFunc("/api/sessions", makeHttpHandleFunc(sessionClient.GetAllSessions)).Methods("GET").Queries("shopid", "{shopid}")
	protected.HandleFunc("/api/active-sessions", makeHttpHandleFunc(sessionClient.GetAllActiveSessions)).Methods("GET").Queries("shopid", "{shopid}")
	protected.HandleFunc("/api/sessions/{sessionid}", makeHttpHandleFunc(sessionClient.GetSessionById)).Methods("GET")
	protected.HandleFunc("/api/sessions/create-join", makeHttpHandleFunc(sessionClient.CreateSessionAdmin)).Methods("POST") //create or join session admin
	protected.HandleFunc("/api/sessions", makeHttpHandleFunc(sessionClient.UpdateSession)).Methods("PUT")

	//ORDERS
	session.HandleFunc("/api/orders", makeHttpHandleFunc(orderClient.GetOrdersBySessionId)).Methods("GET")               //customer
	protected.HandleFunc("/api/orders/{sessionid}", makeHttpHandleFunc(orderClient.GetOrdersBySessionId)).Methods("GET") //shopowner
	protected.HandleFunc("/api/order/{orderid}", makeHttpHandleFunc(orderClient.GetOrderById)).Methods("GET")
	session.HandleFunc("/api/order", makeHttpHandleFunc(orderClient.CreateOrder)).Methods("POST")              //customer
	protected.HandleFunc("/api/order-admin", makeHttpHandleFunc(orderClient.CreateOrderAdmin)).Methods("POST") //shopowner
	// protected.HandleFunc("/api/order/{sessionid}", makeHttpHandleFunc(orderClient.CreateOrder)).Methods("POST") //shopowner
	protected.HandleFunc("/api/order", makeHttpHandleFunc(orderClient.UpdateOrder)).Methods("PUT")
	protected.HandleFunc("/api/order/cancel", makeHttpHandleFunc(orderClient.CancelOrder)).Methods("PUT")

	//PAYMENTS
	protected.HandleFunc("/api/payment", makeHttpHandleFunc(paymentClient.CreatePayment)).Methods("POST")
	session.HandleFunc("/api/payment", makeHttpHandleFunc(paymentClient.GetPayment)).Methods("GET")

	//HISTORY
	protected.HandleFunc("/api/history", makeHttpHandleFunc(historyClient.GetSessionSummary)).Methods("GET").Queries("shopid", "{shopid}", "offset", "{offset}", "limit", "{limit}", "time-range", "{time-range}")

	return http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
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

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if id := r.Context().Value(middleware.CtxKeyID); id != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`"Verification successful"`))
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`"Unauthorized"`))
}
