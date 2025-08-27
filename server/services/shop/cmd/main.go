package main

import (
	"log"
	"time"

	shop "github.com/priyanshubiswas1/shopservice"
	"github.com/priyanshubiswas1/shopservice/pkg"
	"github.com/tinrab/retry"
)

func main() {
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}
	var r shop.IRepository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = shop.NewPostgresRepository(cfg.DbUrl)
		if err != nil {
			log.Println("failed to connect to db: ", err)
			return err
		}
		return
	})
	defer r.Close()
	log.Println("shop repo created")

	service, err := shop.NewService(r)
	if err != nil {
		log.Fatal("failed to create service: ", err)
	}
	log.Println("shop service created")

	log.Println("starting grpc server on port: ", cfg.Port)
	if err := shop.ListenGRPC(service, cfg.Port); err != nil {
		log.Fatal("failed to create grpc server: ", err)
	}
	log.Println("listening on port: ", cfg.Port)
}
