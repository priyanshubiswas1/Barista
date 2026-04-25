package main

import (
	"context"
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

	// RabbitMQ + Garage + Hunyuan3D worker (optional, enabled when configured).
	var scanPublisher *menu.ScanPublisher
	var garageUploader *pkg.GarageUploader

	if cfg.RMQHost != "" && cfg.RMQUsername != "" && cfg.RMQPassword != "" && cfg.RMQVhost != "" {
		log.Println("creating rmq pub connection (3d scan)")
		rmqPubConn, err := menu.ConnectRabbitMQ(cfg.RMQUsername, cfg.RMQPassword, cfg.RMQHost, cfg.RMQVhost)
		if err != nil {
			log.Fatal("failed to connect to rabbitmq server (pub): ", err)
		}
		defer rmqPubConn.Close()

		exchange := cfg.ScanExchange
		if exchange == "" {
			exchange = "photogrammetry_event"
		}
		routingKey := cfg.ScanRoutingKey
		if routingKey == "" {
			routingKey = "hunyuan_3dscan"
		}
		scanPublisher = menu.NewScanPublisher(rmqPubConn, exchange, routingKey)

		if cfg.GarageEndpoint != "" && cfg.GarageAccessKey != "" && cfg.GarageSecretKey != "" && cfg.GarageBucket != "" {
			region := cfg.GarageRegion
			if region == "" {
				region = "garage"
			}
			garageUploader, err = pkg.NewGarageUploader(context.Background(), pkg.GarageConfig{
				Endpoint:        cfg.GarageEndpoint,
				Region:          region,
				AccessKey:       cfg.GarageAccessKey,
				SecretKey:       cfg.GarageSecretKey,
				Bucket:          cfg.GarageBucket,
				PublicURLPrefix: cfg.GaragePublicURLPrefix,
			})
			if err != nil {
				log.Fatal("failed to init garage uploader: ", err)
			}

			log.Println("creating rmq consumer connection (3d scan)")
			rmqConsumeConn, err := menu.ConnectRabbitMQ(cfg.RMQUsername, cfg.RMQPassword, cfg.RMQHost, cfg.RMQVhost)
			if err != nil {
				log.Fatal("failed to connect to rabbitmq server (consume): ", err)
			}
			defer rmqConsumeConn.Close()

			queue := cfg.ScanQueueName
			if queue == "" {
				queue = "hunyuan_3dscan_queue"
			}
			consumerName := cfg.ConsumerName
			if consumerName == "" {
				consumerName = "menu_scan_worker"
			}
			python := cfg.HunyuanPython
			if python == "" {
				python = "python"
			}
			script := cfg.HunyuanScript
			if script == "" {
				script = "/app/services/photogrammetry/hunyuan3d.py"
			}

			if err := menu.StartScanConsumer(r, garageUploader, rmqConsumeConn, exchange, queue, routingKey, consumerName, python, script); err != nil {
				log.Fatal("failed to start scan consumer: ", err)
			}
			log.Println("3d scan worker started")
		} else {
			log.Println("garage not configured - 3d scan uploads will be disabled")
		}
	} else {
		log.Println("rabbitmq not configured - 3d scan jobs disabled")
	}

	service := menu.NewService(r, cloudinaryService, scanPublisher)
	log.Println("menu service created")

	if err = menu.ListenGRPC(service, cfg.Port); err != nil {
		log.Fatal("failed to start server: ", err)
	}
	log.Println("Listening to port: ", cfg.Port)
}
