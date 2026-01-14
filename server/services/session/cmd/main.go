package main

import (
	"log"
	"time"

	session "github.com/priyanshubiswas1/barista/sessionservice"
	"github.com/priyanshubiswas1/barista/sessionservice/client"
	"github.com/priyanshubiswas1/barista/sessionservice/pkg"
	"github.com/tinrab/retry"
)

func main() {
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}
	// Environment Variable Validation
	validateEnv(&cfg)

	var r session.IRepository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = session.NewPostgresRepository(cfg.DbUrl)
		if err != nil {
			log.Println("failed to connect to db: ", err)
			return err
		}
		return
	})
	defer r.Close()
	log.Println("session repo created")

	log.Println("starting redis client")
	redisClient, err := session.NewRedisClient(cfg.RedisUrl)
	if err != nil {
		log.Fatal("failed to create redis client: ", err)
	}
	log.Println("redis client created")

	log.Println("creating table grpc client")
	tableClient, err := client.NewTable(cfg.TableService)
	if err != nil {
		log.Printf("Warning: failed to create table grpc client: %v", err)
	} else {
		log.Println("table grpc client created")
	}

	log.Println("creating history grpc client")
	historyClient, err := client.NewHistory(cfg.HistoryService)
	if err != nil {
		log.Printf("Warning: failed to create history grpc client: %v", err)
	} else {
		log.Println("history grpc client created")
	}

	service, err := session.NewService(r, redisClient, tableClient, historyClient)
	if err != nil {
		log.Fatal("failed to create service: ", err)
	}
	log.Println("session service created")

	log.Println("starting grpc server on port: ", cfg.Port)
	if err := session.ListenGRPC(service, cfg.Port); err != nil {
		log.Fatal("failed to create grpc server: ", err)
	}
	log.Println("listening on port: ", cfg.Port)
}

func validateEnv(cfg *pkg.Config) {
	required := map[string]string{
		"PORT":                cfg.Port,
		"DB_URL":              cfg.DbUrl,
		"REDIS_URL":           cfg.RedisUrl,
		"TABLE_SERVICE_URL":   cfg.TableService,
		"HISTORY_SERVICE_URL": cfg.HistoryService,
	}

	for name, val := range required {
		if val == "" {
			log.Fatalf("Environment variable %s is required", name)
		}
	}
}
