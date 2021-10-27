package queue

import (
	"encoding/json"
	"fmt"

	"github.com/Konstantsiy/image-converter/internal/config"
	"github.com/streadway/amqp"
)

// MIMEContentType represents the content type of publishing message.
const MIMEContentType = "application/json"

// RabbitMQProducer implements RabbitMQ queue producer.
type RabbitMQProducer struct {
	client *rabbitMQClient
}

// NewRabbitMQProducer creates new RabbitMQ producer.
func NewRabbitMQProducer(conf *config.RabbitMQConfig) (*RabbitMQProducer, error) {
	client, err := initRabbitMQClient(conf)
	if err != nil {
		return nil, err
	}

	return &RabbitMQProducer{client: client}, nil
}

// SendToQueue sends messages to the queue.
func (p *RabbitMQProducer) SendToQueue(fileID, filename, sourceFormat, targetFormat, requestID string, ratio int) error {
	msg := queueMessage{
		FileID:       fileID,
		Filename:     filename,
		SourceFormat: sourceFormat,
		TargetFormat: targetFormat,
		RequestID:    requestID,
		Ratio:        ratio,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("can't marshal queue message: %w", err)
	}

	err = p.client.ch.Publish("", p.client.queue.Name, false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  MIMEContentType,
			Body:         body,
		})
	if err != nil {
		return fmt.Errorf("can't publish queue message: %w", err)
	}

	return nil
}
