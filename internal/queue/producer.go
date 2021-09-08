package queue

import (
	"encoding/json"
	"fmt"

	"github.com/Konstantsiy/image-converter/internal/config"
	"github.com/streadway/amqp"
)

// Publisher
type Publisher struct {
	client *rabbitMQClient
}

func NewPublisher(conf *config.Config) (*Publisher, error) {
	client, err := initRabbitMQClient(conf)
	if err != nil {
		return nil, err
	}

	return &Publisher{client: client}, nil
}

// SendToQueue sends messages to the queue.
func (p *Publisher) SendToQueue(fileID, filename, sourceFormat, targetFormat, requestID string, ratio int) error {
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
			ContentType:  "application/json",
			Body:         body,
		})
	if err != nil {
		return fmt.Errorf("can't publish queue message: %w", err)
	}

	return nil
}
