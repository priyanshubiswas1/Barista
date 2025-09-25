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

	server, err := NewAPIServer(cfg.AuthClient, cfg.ShopClient, cfg.MenuClient, cfg.TableClient, cfg.SessionClient, cfg.OrderClient)
	if err != nil {
		log.Fatal("failed to create server")
	}
	log.Println("server created: ", server.authClientUrl, server.shopClientUrl, server.menuClientUrl, server.tableClientUrl)

	if err = server.ListenServer(cfg.Port); err != nil {
		log.Fatal("error while listening to port : ", err)
	}
	log.Println("listening to port: ", cfg.Port)
}
