package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	auth "github.com/priyanshubiswas1/barista/apigateway/handler"
	"github.com/priyanshubiswas1/barista/apigateway/middleware"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

type APIServer struct {
	authClientUrl string
	shopClientUrl string
}

func NewAPIServer(authClient, shopClientUrl string) (*APIServer, error) {
	return &APIServer{
		authClientUrl: authClient,
		shopClientUrl: shopClientUrl,
	}, nil
}

func (s *APIServer) ListenServer(port string) error {
	authClient, err := auth.NewClient(s.authClientUrl)
	// shopClient, err := auth.NewClient(s.shopClientUrl)
	if err != nil {
		return err
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/auth/register", makeHttpHandleFunc(authClient.RegisterHandler)).Methods("POST")
	router.HandleFunc("/api/auth/login", makeHttpHandleFunc(authClient.LoginHandler)).Methods("POST")

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/api/auth/user", makeHttpHandleFunc(authClient.GetUserHandler)).Methods("GET")
	protected.HandleFunc("/api/auth/user", makeHttpHandleFunc(authClient.UpdateUserHandler)).Methods("PUT")
	protected.HandleFunc("/api/auth/user", makeHttpHandleFunc(authClient.DeleteUserHandler)).Methods("DELETE")
	protected.HandleFunc("/api/shop", makeHttpHandleFunc(authClient.GetUserHandler)).Methods("GET")
	protected.HandleFunc("/api/shops", makeHttpHandleFunc(authClient.GetUserHandler)).Methods("GET")
	protected.HandleFunc("/api/shop", makeHttpHandleFunc(authClient.UpdateUserHandler)).Methods("POST")
	protected.HandleFunc("/api/shop", makeHttpHandleFunc(authClient.UpdateUserHandler)).Methods("PUT")
	protected.HandleFunc("/api/shop", makeHttpHandleFunc(authClient.DeleteUserHandler)).Methods("DELETE")

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
