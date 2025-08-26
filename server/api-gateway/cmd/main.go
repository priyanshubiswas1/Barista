package main

import (
	"log"

	"github.com/priyanshubiswas1/barista/apigateway/pkg"
)

func main() {
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}
	log.Println("config loaded")

	server, err := NewAPIServer(cfg.AuthClient, cfg.ShopClient)
	if err != nil {
		log.Fatal("failed to create server")
	}
	log.Println("server created: ", cfg.AuthClient, cfg.ShopClient)

	if err = server.ListenServer(cfg.Port); err != nil {
		log.Fatal("error while listening to port : ", err)
	}
	log.Println("listening to port: ", cfg.Port)
}
