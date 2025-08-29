package main

import (
	"log"
	"time"

	auth "github.com/priyanshubiswas1/barista/authservice"
	"github.com/priyanshubiswas1/barista/authservice/pkg"
	"github.com/tinrab/retry"
)

func main() {
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}
	var r auth.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = auth.NewPostgresRepository(cfg.DbUrl)
		if err != nil {
			log.Println("failed to connect to db: ", err)
			return err
		}
		return
	})
	defer r.Close()
	log.Println("auth repo created")

	s := auth.NewService(r)
	log.Println("starting grpc server on port: ", cfg.Port)
	if err := auth.ListenGRPC(s, cfg.Port); err != nil {
		log.Fatal("failed to create grpc server: ", err)
	}
	log.Println("listening on port: ", cfg.Port)
}
