package main

import (
	"log"
	"time"

	menu "github.com/priyanshubiswas1/barista/menuservice"
	"github.com/priyanshubiswas1/barista/menuservice/pkg"
	"github.com/tinrab/retry"
)

func main() {
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}

	var r menu.IRepository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = menu.NewPostgresRepository(cfg.DbUrl)
		if err != nil {
			log.Println("failed to connect to db: ", err)
			return err
		}
		return
	})
	defer r.Close()
	log.Println("menu repo created")

	var cloudinaryService *pkg.CloudinaryService
	if cfg.CloudinaryApiKey != "" && cfg.CloudinarySecret != "" && cfg.CloudinaryCloudName != "" {
		cloudinaryService = pkg.NewCloudinaryService(cfg.CloudinaryApiKey, cfg.CloudinarySecret, cfg.CloudinaryCloudName)
		log.Println("cloudinary service initialized")
	} else {
		log.Println("cloudinary service not configured - image uploads will fail")
	}

	service := menu.NewService(r, cloudinaryService)
	log.Println("menu service created")

	if err = menu.ListenHTTP(service, cfg.Port); err != nil {
		log.Fatal("failed to start server: ", err)
	}
	log.Println("Listening to port: ", cfg.Port)
}
