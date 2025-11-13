package main

import (
	"log"
	"time"

	order "github.com/priyanshubiswas1/barista/orderservice"
	"github.com/priyanshubiswas1/barista/orderservice/client"
	"github.com/priyanshubiswas1/barista/orderservice/pkg"
	"github.com/tinrab/retry"
)

func main() {
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}
	var r order.IRepository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = order.NewPostgresRepository(cfg.DbUrl)
		if err != nil {
			log.Println("failed to connect to db: ", err)
			return err
		}
		return
	})
	defer r.Close()
	log.Println("order repo created")

	log.Println("creating session grpc client")
	sessionClient, err := client.NewSession(cfg.SessionClient)
	if err != nil {
		log.Fatal("failed to create session grpc client: ", err)
	}
	log.Println("session grpc client created")

	log.Println("creating history grpc client")
	historyClient, err := client.NewHistory(cfg.HistoryClient)
	if err != nil {
		log.Fatal("failed to create history grpc client: ", err)
	}
	log.Println("history grpc client created")

	service, err := order.NewService(r, sessionClient, historyClient, cfg.WebSocketService)
	if err != nil {
		log.Fatal("failed to create service: ", err)
	}
	log.Println("order service created")

	log.Println("starting grpc server on port: ", cfg.Port)
	if err := order.ListenGRPC(service, cfg.Port); err != nil {
		log.Fatal("failed to create grpc server: ", err)
	}
	log.Println("listening on port: ", cfg.Port)

}
