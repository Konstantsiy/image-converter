// Package queue implements the message queue functionality using RabbitMQ.
package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/Konstantsiy/image-converter/internal/config"
	"github.com/Konstantsiy/image-converter/internal/converter"
	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/internal/storage"
	"github.com/Konstantsiy/image-converter/pkg/logger"
	"github.com/streadway/amqp"
)

// Consumer listens to the queue and processes outgoing messages.
type Consumer struct {
	client  *rabbitMQClient
	repo    *repository.Repository
	storage *storage.Storage
}

func NewConsumer(repo *repository.Repository, storage *storage.Storage, conf *config.RabbitMQConfig) (*Consumer, error) {
	client, err := initRabbitMQClient(conf)
	if err != nil {
		return nil, err
	}
	return &Consumer{repo: repo, storage: storage, client: client}, nil
}

// LaunchListener listens to the queue channel in a separate goroutine.
func (c *Consumer) LaunchListener() error {
	err := c.client.ch.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("can't configure QoS: %w", err)
	}

	msgChannel, err := c.client.ch.Consume(c.client.queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("can't register channel: %w", err)
	}

	for {
		msg := <-msgChannel
		go func() {
			err := c.consumeFromQueue(&msg)
			if err != nil {
				nErr := msg.Nack(true, false)
				if nErr != nil {
					logger.Error(context.Background(), fmt.Errorf("can't make negative acknowledgement: %w, (original error: %v)", nErr, err))
				}
				return
			}
			aErr := msg.Ack(false)
			if aErr != nil {
				logger.Error(context.Background(), fmt.Errorf("can't make acknowledgement: %w", aErr))
			}
		}()
	}
}

// consumeFromQueue wraps the message processing and confirms its completion.
func (c *Consumer) consumeFromQueue(msg *amqp.Delivery) error {
	var data queueMessage
	err := json.NewDecoder(bytes.NewReader(msg.Body)).Decode(&data)
	if err != nil {
		return fmt.Errorf("can't decode queue message: %w", err)
	}

	err = c.process(data)
	if err != nil {
		uErr := c.repo.UpdateRequest(data.RequestID, repository.RequestStatusFailed, "")
		if uErr != nil {
			logger.Error(context.Background(), fmt.Errorf("can't update request with id %s: %w, (original error: %v)", data.RequestID, err, err))
		}

		return fmt.Errorf("error processing a message from the queue: %w", err)
	}

	return nil
}

// process processes the current message from the queue.
func (c *Consumer) process(data queueMessage) error {
	sourceFile, err := c.storage.DownloadFile(data.FileID)
	if err != nil {
		return fmt.Errorf("storage error: %w", err)
	}

	targetFile, err := converter.Convert(sourceFile, data.TargetFormat, data.Ratio)
	if err != nil {
		return fmt.Errorf("converter error: %w", err)
	}

	err = c.repo.UpdateRequest(data.RequestID, repository.RequestStatusProcessing, "")
	if err != nil {
		return fmt.Errorf("can't update request with id %s: %w", data.RequestID, err)
	}

	targetFileID, err := c.repo.InsertImage(data.Filename, data.TargetFormat)
	if err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	err = c.storage.UploadFile(targetFile, targetFileID)
	if err != nil {
		return fmt.Errorf("storage error: %w", err)
	}

	err = c.repo.UpdateRequest(data.RequestID, repository.RequestStatusDone, targetFileID)
	if err != nil {
		return fmt.Errorf("request updating error: %w", err)
	}

	return nil
}
