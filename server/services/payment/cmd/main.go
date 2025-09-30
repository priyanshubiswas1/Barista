package main

import (
	"log"
	"time"

	payment "github.com/priyanshubiswas1/barista/paymentservice"
	"github.com/priyanshubiswas1/barista/paymentservice/client"
	"github.com/priyanshubiswas1/barista/paymentservice/pkg"
	"github.com/tinrab/retry"
)

func main() {
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}
	var r payment.IRepository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = payment.NewPostgresRepository(cfg.DbUrl)
		if err != nil {
			log.Println("failed to connect to db: ", err)
			return err
		}
		return
	})
	defer r.Close()
	log.Println("payment repo created")

	log.Println("creating session grpc client")
	sessionClient, err := client.NewSession(cfg.SessionClient)
	if err != nil {
		log.Fatal("failed to create session grpc client: ", err)
	}
	log.Println("session grpc client created")

	service, err := payment.NewService(r, sessionClient)
	if err != nil {
		log.Fatal("failed to create service: ", err)
	}
	log.Println("payment service created")

	log.Println("starting grpc server on port: ", cfg.Port)
	if err := payment.ListenGRPC(service, cfg.Port); err != nil {
		log.Fatal("failed to create grpc server: ", err)
	}
	log.Println("listening on port: ", cfg.Port)

}
