package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/priyanshubiswas1/barista/apigateway/pkg"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to get env")
	}
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config")
	}

	server, err := NewAPIServer(cfg.AuthClient)
	if err != nil {
		log.Fatal("failed to create server")
	}

	if err = server.ListenServer(cfg.Port); err != nil {
		log.Fatal("error while listening to port : ", err)
	}
	log.Println("listening to port: ", cfg.Port)
}
