package menu

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

