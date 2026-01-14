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

	// Environment Variable Validation
	validateEnv(&cfg)
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
	sessionClient, err := client.NewSession(cfg.SessionService)
	if err != nil {
		log.Printf("Warning: failed to create session grpc client: %v", err)
	} else {
		log.Println("session grpc client created")
	}

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

func validateEnv(cfg *pkg.Config) {
	required := map[string]string{
		"PORT":                cfg.Port,
		"DB_URL":              cfg.DbUrl,
		"SESSION_SERVICE_URL": cfg.SessionService,
	}

	for name, val := range required {
		if val == "" {
			log.Fatalf("Environment variable %s is required", name)
		}
	}
}
