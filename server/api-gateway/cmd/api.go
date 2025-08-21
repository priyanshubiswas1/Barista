package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	auth "github.com/priyanshubiswas1/barista/apigateway/handler"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

type APIServer struct {
	authClientUrl string
}

func NewAPIServer(authClient string) (*APIServer, error) {
	return &APIServer{
		authClientUrl: authClient,
	}, nil
}

func (s *APIServer) ListenServer(port string) error {
	authClient, err := auth.NewClient(s.authClientUrl)
	if err != nil {
		return err
	}

	router := mux.NewRouter()
	router.HandleFunc("/register", makeHttpHandleFunc(authClient.RegisterHandler)).Methods("POST")
	router.HandleFunc("/login", makeHttpHandleFunc(authClient.LoginHandler)).Methods("POST")
	router.HandleFunc("/user", makeHttpHandleFunc(authClient.GetUserHandler)).Methods("GET")
	router.HandleFunc("/user", makeHttpHandleFunc(authClient.UpdateUserHandler)).Methods("PUT")
	router.HandleFunc("/user", makeHttpHandleFunc(authClient.DeleteUserHandler)).Methods("DELETE")

	return http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}

func makeHttpHandleFunc(apiFunc APIFunc) http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		if err := apiFunc(ctx, w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		}
	}
}

func writeJSON(w http.ResponseWriter, s int, data any) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(data)
}
