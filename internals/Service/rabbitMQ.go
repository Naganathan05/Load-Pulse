package Service

import (
	"fmt"
	"log"
	"sync"
	config "Load-Pulse/Config"
	"github.com/streadway/amqp"
)

var (
	connection *amqp.Connection
	once       sync.Once
)

func ConnectRabbitMQ() {
	var err error;
	cfg := config.GetConfig();
	once.Do(func() {
		LogServer("[LOG]: Establishing RabbitMQ Connection\n");
		connection, err = amqp.Dial(cfg.RabbitMQURL);
		if err != nil {
			log.Fatalf("[ERR]: Failed to connect to RabbitMQ: %s", err);
		}
		LogServer("[LOG]: RabbitMQ Connection Established\n");
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

	logMsg := fmt.Sprintf("[LOG]: Published Stats Events to %s Succussfully.\n", queueName);
	LogServer(logMsg);
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
	return nil;
}

func InspectQueue(queueName string) (amqp.Queue, error) {
	channel, err := connection.Channel()
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("[ERR]: Failed to open channel: %v", err)
	}
	defer channel.Close()

	queue, err := channel.QueueInspect(queueName)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("[ERR]: Failed to inspect queue: %v", err)
	}

	return queue, nil
}

func PublishToQueue(queueName string, message []byte) error {

	if err := CreateQueue(queueName); err != nil {
        return err
    }
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
	return nil;
}

func ConsumeFromQueue(queueName string) (<-chan amqp.Delivery, error) {
	if err := CreateQueue(queueName); err != nil {
        return nil, err
    }
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
		LogServer("[LOG]: RabbitMQ Connection Closed.\n");
	}
}