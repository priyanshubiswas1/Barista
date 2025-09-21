package main

import (
	"log"
	"sync"

	internal "github.com/priyanshubiswas1/barista/qrgenservice"
	"github.com/priyanshubiswas1/barista/qrgenservice/pkg"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}

	// PUBLISH CONN
	log.Println("starting pub rabbitmq connection")
	rmqPubConn, err := internal.ConnectRabbitMQ(cfg.RMQUsername, cfg.RMQPassword, cfg.RMQHost, cfg.RMQVhost)
	if err != nil {
		log.Fatal("failed to connect to rabbitmq server: ", err)
	}
	defer rmqPubConn.Close()

	rmqPubClient, err := internal.NewRabbitMQClient(rmqPubConn)
	if err != nil {
		log.Fatal("failed to create rabbitmq client: ", err)
	}
	defer rmqPubClient.Close()

	// CONSUME
	log.Println("starting con rabbitmq connection")
	rmqConn, err := internal.ConnectRabbitMQ(cfg.RMQUsername, cfg.RMQPassword, cfg.RMQHost, cfg.RMQVhost)
	if err != nil {
		log.Fatal("failed to connect to rabbitmq server: ", err)
	}
	defer rmqConn.Close()

	rmqClient, err := internal.NewRabbitMQClient(rmqConn)
	if err != nil {
		log.Fatal("failed to create rabbitmq client: ", err)
	}
	defer rmqClient.Close()

	if err := rmqClient.CreateQueue(cfg.ConsumeQueueName, true, false); err != nil {
		log.Fatal("failed to create rabbitmq queue: ", err)
	}

	if err := rmqClient.CreateBinding(cfg.ConsumeQueueName, cfg.ConsumeBindingKey, cfg.Exchange); err != nil {
		log.Fatal("failed to create rabbitmq binding: ", err)
	}

	// Start the consumer
	messageBus, err := rmqClient.StartConsumer(cfg.ConsumeQueueName, cfg.ConsumerName, false)
	if err != nil {
		log.Fatal("failed to start consumer: ", err)
	}

	var wg sync.WaitGroup

	go func() {
		for message := range messageBus {
			log.Println("Waiting for messages...")

			wg.Add(1)
			go func(msg amqp091.Delivery) {
				defer wg.Done()
				log.Println("Received message:", msg.CorrelationId)

				// Process the message (acknowledge and send response)
				if err := msg.Ack(false); err != nil {
					log.Fatal("failed to acknowledge message: ", err)
					return
				}

				//Generate QR
				qrCode := internal.GenerateQR(string(msg.Body))

				// Send the message back after processing
				if err := rmqPubClient.Send("gen_qr_event", msg.ReplyTo, amqp091.Publishing{
					ContentType:   "text/plain",
					DeliveryMode:  amqp091.Transient,
					Body:          qrCode,
					CorrelationId: msg.CorrelationId,
				}); err != nil {
					log.Fatal("failed to send generated QR: ", err)
					return
				}

				log.Println("Message processed and acknowledged:", msg.MessageId)
			}(message)
		}
	}()

	log.Println("Service is running, press CTRL+C to exit.")
	select {}
}
