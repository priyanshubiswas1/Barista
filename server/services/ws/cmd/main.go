package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/priyanshubiswas1/barista/ws"
	"github.com/priyanshubiswas1/barista/ws/pkg"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

func main() {
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}
	log.Println("config loaded")

	manager := ws.NewManager(cfg.AllowedOrigins)
	http.HandleFunc("/ws", manager.ServeWS)
	http.HandleFunc("/send-event", makeHttpHandleFunc(manager.RouteEvent))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))
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
