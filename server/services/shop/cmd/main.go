package main

import (
	"log"
	"time"

	shop "github.com/priyanshubiswas1/barista/shopservice"
	"github.com/priyanshubiswas1/barista/shopservice/pkg"
	"github.com/tinrab/retry"
)

func main() {
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}
	// Environment Variable Validation
	validateEnv(cfg)

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

	log.Println("starting http server on port: ", cfg.Port)
	if err := shop.ListenHTTP(service, cfg.Port); err != nil {
		log.Fatal("failed to create http server: ", err)
	}
	log.Println("listening on port: ", cfg.Port)
}

func validateEnv(cfg *pkg.Config) {
	required := map[string]string{
		"PORT":   cfg.Port,
		"DB_URL": cfg.DbUrl,
	}

	for name, val := range required {
		if val == "" {
			log.Fatalf("Environment variable %s is required", name)
		}
	}
}
