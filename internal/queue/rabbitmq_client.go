package queue

import (
	"fmt"

	"github.com/Konstantsiy/image-converter/internal/config"
	"github.com/streadway/amqp"
)

// queueMessage represents a message that is passed to the queue.
type queueMessage struct {
	FileID       string
	Filename     string
	SourceFormat string
	TargetFormat string
	RequestID    string
	Ratio        int
}

// rabbitMQClient provides connection to the queue via a specific channel.
type rabbitMQClient struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue *amqp.Queue
}

// initRabbitMQClient initializes the queue client.
func initRabbitMQClient(conf *config.RabbitMQConfig) (*rabbitMQClient, error) {
	if conf.AMQPConnectionURL == "" || conf.QueueName == "" {
		return nil, fmt.Errorf("RabbitMQ configurations should not be empty")
	}

	conn, err := amqp.Dial(conf.AMQPConnectionURL)
	if err != nil {
		return nil, fmt.Errorf("can't connect to AMQP: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("can't create an AMQP channel: %w", err)
	}

	queue, err := ch.QueueDeclare(conf.QueueName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("can't declare %s queue", conf.QueueName)
	}

	return &rabbitMQClient{conn: conn, ch: ch, queue: &queue}, nil
}
