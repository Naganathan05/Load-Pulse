package Service

import (
	"fmt"
	"sync"
	"github.com/streadway/amqp"
)

var (
	connection *amqp.Connection
	channel    *amqp.Channel
	once       sync.Once
)

func ConnectRabbitMQ() error {
	var err error;
	once.Do(func() {
		fmt.Println("[LOG]: Establishing RabbitMQ Connection");
		connection, err = amqp.Dial("amqp://guest:guest@localhost:5672/");
		if err != nil {
			err = fmt.Errorf("[ERR]: Failed to Connect to RabbitMQ: %s", err);
			return;
		}

		channel, err = connection.Channel();
		if err != nil {
			err = fmt.Errorf("[ERR]: Failed to Open a Channel: %s", err);
			return;
		}
		fmt.Println("[LOG]: RabbitMQ Connection Established.");
	});

	return nil;
}

func CreateQueue(queueName string) error {
	_, err := channel.QueueDeclare(
		queueName, // queue name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("[ERR]: Failed to Declare Queue: %v", err);
	}
	return nil;
}

func PublishToQueue(queueName string, message []byte) error {
	err := channel.Publish(
		"",          // exchange
		queueName,   // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("[ERR]: Failed to Publish Message: %v", err);
	}
	return nil;
}

func ConsumeFromQueue(queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, fmt.Errorf("[ERR]: Failed to Consume Messages: %v", err);
	}
	return msgs, nil;
}

func CloseRabbitMQ() {
	if channel != nil {
		channel.Close();
		fmt.Println("[LOG]: RabbitMQ Channel Closed.");
	}
	if connection != nil {
		connection.Close();
		fmt.Println("[LOG]: RabbitMQ Connection Closed.");
	}
}