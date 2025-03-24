package Service

import (
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

var (
	connection *amqp.Connection
	once       sync.Once
)

func ConnectRabbitMQ() {
	var err error;
	once.Do(func() {
		fmt.Println("[LOG]: Establishing RabbitMQ Connection");
		connection, err = amqp.Dial("amqp://guest:guest@localhost:5672/");
		if err != nil {
			log.Fatalf("[ERR]: Failed to connect to RabbitMQ: %s", err);
		}
		fmt.Println("[LOG]: RabbitMQ Connection Established");
	})
}

func CreateQueue(queueName string) error {
	channel, err := connection.Channel();
	if err != nil {
		return fmt.Errorf("[ERR]: Failed to open channel: %v", err);
	}
	defer channel.Close();

	_, err = channel.QueueDeclare(
		queueName, // queue name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("[ERR]: Failed to declare queue: %v", err);
	}

	fmt.Printf("[LOG]: Queue %s is ready.\n", queueName);
	return nil;
}

func DeleteQueue(queueName string) error {
	channel, err := connection.Channel();
	if err != nil {
		return fmt.Errorf("[ERR]: Failed to open channel: %v", err);
	}
	defer channel.Close();

	_, err = channel.QueueDelete(
		queueName, // queue name
		false,     // ifUnused
		false,     // ifEmpty
		false,     // noWait
	)
	if err != nil {
		return fmt.Errorf("[ERR]: Failed to delete queue: %v", err);
	}

	fmt.Printf("[LOG]: Queue %s deleted successfully.\n", queueName);
	return nil;
}

// PublishToQueue sends a message to the specified queue.
func PublishToQueue(queueName string, message []byte) error {
	channel, err := connection.Channel();
	if err != nil {
		return fmt.Errorf("[ERR]: Failed to open channel: %v", err);
	}
	defer channel.Close();

	err = channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("[ERR]: Failed to publish message: %v", err);
	}
	return nil
}

func ConsumeFromQueue(queueName string) (<-chan amqp.Delivery, error) {
	channel, err := connection.Channel();
	if err != nil {
		return nil, fmt.Errorf("[ERR]: Failed to open channel: %v", err);
	}

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
		channel.Close();
		return nil, fmt.Errorf("[ERR]: Failed to start consuming messages: %v", err);
	}

	return msgs, nil;
}

func CloseRabbitMQ() {
	if connection != nil {
		connection.Close();
		fmt.Println("[LOG]: RabbitMQ Connection Closed.");
	}
}