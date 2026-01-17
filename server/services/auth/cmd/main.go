package main

import (
	"log"
	"time"

	auth "github.com/priyanshubiswas1/barista/authservice"
	"github.com/priyanshubiswas1/barista/authservice/client"
	"github.com/priyanshubiswas1/barista/authservice/pkg"
	"github.com/tinrab/retry"
)

func main() {
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}
	// Environment Variable Validation
	validateEnv(cfg)

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

	// auth service is useless without shop service and menu service
	shopClient, err := client.NewShop(cfg.ShopService)
	if err != nil {
		log.Printf("Warning: failed to connect to shop service: %v", err)
	}

	menuClient, err := client.NewMenu(cfg.MenuService)
	if err != nil {
		log.Printf("Warning: failed to connect to menu service: %v", err)
	}

	s := auth.NewService(r, shopClient, menuClient)
	log.Println("starting http server on port: ", cfg.Port)
	if err := auth.ListenHTTP(s, cfg.Port); err != nil {
		log.Fatal("failed to create http server: ", err)
	}
	log.Println("listening on port: ", cfg.Port)
}

func validateEnv(cfg *pkg.Config) {
	required := map[string]string{
		"PORT":             cfg.Port,
		"DB_URL":           cfg.DbUrl,
		"SHOP_SERVICE_URL": cfg.ShopService,
		"MENU_SERVICE_URL": cfg.MenuService,
	}

	for name, val := range required {
		if val == "" {
			log.Fatalf("Environment variable %s is required", name)
		}
	}
}
