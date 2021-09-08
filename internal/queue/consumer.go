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

func NewConsumer(repo *repository.Repository, storage *storage.Storage, conf *config.Config) (*Consumer, error) {
	client, err := initRabbitMQClient(conf)
	logger.Info(context.Background(), "consumer client is init")
	if err != nil {
		return nil, err
	}
	return &Consumer{repo: repo, storage: storage, client: client}, nil
}

// ListenToQueue listens to the queue channel in a separate goroutine.
func (c *Consumer) ListenToQueue() error {
	err := c.client.ch.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("can't configure QoS: %w", err)
	}

	logger.Info(context.Background(), "QoS is configured")

	msgChannel, err := c.client.ch.Consume(c.client.queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("can't register channel: %w", err)
	}

	logger.Info(context.Background(), "queue channel is created")

	go func() {
		for {
			msg := <-msgChannel
			go c.processQueueMessage(&msg)
		}
	}()

	return nil
}

// processQueueMessage processes the current message from the queue.
func (c *Consumer) processQueueMessage(msg *amqp.Delivery) {
	var data queueMessage
	err := json.NewDecoder(bytes.NewReader(msg.Body)).Decode(&data)
	if err != nil {
		logger.Error(context.Background(), fmt.Errorf("can't decode queue message: %w", err))
		return
	}

	err = msg.Ack(false)
	if err != nil {
		logger.Error(context.Background(), err)
	}

	logger.Info(context.Background(), fmt.Sprintf("message: %+v", data))

	requestCompletion := false
	defer func() {
		if !requestCompletion {
			uErr := c.repo.UpdateRequest(data.RequestID, repository.RequestStatusFailed, "")
			if uErr != nil {
				logger.Error(context.Background(), fmt.Errorf("can't update request with id %s: %w", data.RequestID, err))
			}
		}
	}()

	sourceFile, err := c.storage.DownloadFile(data.FileID)
	if err != nil {
		logger.Error(context.Background(), fmt.Errorf("storage error: %v", err))
		return
	}

	targetFile, err := converter.Convert(sourceFile, data.TargetFormat, data.Ratio)
	if err != nil {
		logger.Error(context.Background(), fmt.Errorf("converter error: %w", err))
		return
	}

	err = c.repo.UpdateRequest(data.RequestID, repository.RequestStatusProcessing, "")
	if err != nil {
		logger.Error(context.Background(), fmt.Errorf("can't update request with id %s: %w", data.RequestID, err))
		return
	}

	logger.Info(context.Background(), "request status updated to the Processing")

	targetFileID, err := c.repo.InsertImage(data.Filename, data.TargetFormat)
	if err != nil {
		logger.Error(context.Background(), fmt.Errorf("repository error: %w", err))
		return
	}

	logger.Info(context.Background(), "targetFile ---> repository")

	err = c.storage.UploadFile(targetFile, targetFileID)
	if err != nil {
		logger.Error(context.Background(), fmt.Errorf("storage error: %w", err))
		return
	}

	logger.Info(context.Background(), "targetFile ---> storage")

	err = c.repo.UpdateRequest(data.RequestID, repository.RequestStatusDone, targetFileID)
	if err != nil {
		logger.Error(context.Background(), fmt.Errorf("request updating error: %w", err))
		return
	}

	logger.Info(context.Background(), "request status updated to the Done")

	requestCompletion = true
}
