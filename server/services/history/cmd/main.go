package main

import (
	"log"
	"time"

	history "github.com/priyanshubiswas1/barista/historyservice"
	"github.com/priyanshubiswas1/barista/historyservice/pkg"
	"github.com/tinrab/retry"
)

func main() {
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}

	var r history.IRepository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = history.NewPostgresRepository(cfg.DbUrl)
		if err != nil {
			log.Println("failed to connect to db: ", err)
			return err
		}
		return
	})
	defer r.Close()
	log.Println("history repo created")

	service, err := history.NewService(r)
	if err != nil {
		log.Fatal("failed to create service: ", err)
	}
	log.Println("history service created")

	log.Println("starting grpc server on port: ", cfg.Port)
	if err := history.ListenGRPC(service, cfg.Port); err != nil {
		log.Fatal("failed to create grpc server: ", err)
	}
	log.Println("listening on port: ", cfg.Port)

}
