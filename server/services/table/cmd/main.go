package main

import (
	"log"
	"time"

	table "github.com/priyanshubiswas1/barista/tableservice"
	"github.com/priyanshubiswas1/barista/tableservice/pkg"
	"github.com/tinrab/retry"
)

func main() {
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}

	var r table.IRepository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = table.NewPostgresRepository(cfg.DbUrl)
		if err != nil {
			log.Println("failed to connect to db: ", err)
			return err
		}
		return
	})
	defer r.Close()
	log.Println("table repo created")

	log.Println("creating rmq pub connection")
	rmqPubConn, err := table.ConnectRabbitMQ(cfg.RMQUsername, cfg.RMQPassword, cfg.RMQHost, cfg.RMQVhost)
	if err != nil {
		log.Fatal("failed to connect to rabbitmq server: ", err)
	}
	defer rmqPubConn.Close()

	log.Println("creating rmq consumer connection")
	rmqConsumeConn, err := table.ConnectRabbitMQ(cfg.RMQUsername, cfg.RMQPassword, cfg.RMQHost, cfg.RMQVhost)
	if err != nil {
		log.Fatal("failed to connect to rabbitmq server: ", err)
	}
	defer rmqConsumeConn.Close()

	service := table.NewService(r, rmqPubConn, rmqConsumeConn, cfg.ClientUrl, cfg.Exchange, cfg.ConsumerName)
	log.Println("table service created")

	log.Println("starting grpc server on port: ", cfg.Port)
	if err := table.ListenGRPC(service, cfg.Port); err != nil {
		log.Fatal("failed to create grpc server: ", err)
	}

	log.Println("listening on port: ", cfg.Port)
}
