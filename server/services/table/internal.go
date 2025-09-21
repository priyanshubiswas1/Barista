package table

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, nil
	}
	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rmq *RabbitClient) Close() error {
	return rmq.ch.Close()
}

func (rmq *RabbitClient) CreateQueue(queueName string, durable, autodelete bool) (amqp.Queue, error) {
	queue, err := rmq.ch.QueueDeclare(queueName, durable, autodelete, false, false, amqp.Table{})
	return queue, err
}

func (rmq *RabbitClient) DeleteQueue(queueName string, ifUnused, ifEmpty bool) error {
	_, err := rmq.ch.QueueDelete(queueName, ifUnused, ifEmpty, false)
	return err
}

// name = name of the queue that we want to bind
// binding = binding key
// exchange = exchange name
func (rmq *RabbitClient) CreateBinding(name, bindingkey, exchange string) error {
	return rmq.ch.QueueBind(name, bindingkey, exchange, false, nil)
}

func (rmq *RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	return rmq.ch.PublishWithContext(ctx, exchange, routingKey, true, false, options)
}

func (rmq *RabbitClient) StartConsumer(queueName, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rmq.ch.Consume(queueName, consumer, autoAck, false, false, false, nil)
}
