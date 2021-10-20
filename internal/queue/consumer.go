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

// RabbitMQConsumer listens to the queue and processes outgoing messages.
type RabbitMQConsumer struct {
	client       *rabbitMQClient
	requestsRepo *repository.RequestsRepository
	imagesRepo   *repository.ImagesRepository
	storage      *storage.Storage
}

func NewConsumer(requestsRepo *repository.RequestsRepository, imagesRepo *repository.ImagesRepository, storage *storage.Storage, conf *config.RabbitMQConfig) (*RabbitMQConsumer, error) {
	client, err := initRabbitMQClient(conf)
	if err != nil {
		return nil, err
	}
	return &RabbitMQConsumer{requestsRepo: requestsRepo, imagesRepo: imagesRepo, storage: storage, client: client}, nil
}

// Listen listens to the queue channel in a separate goroutine.
func (c *RabbitMQConsumer) Listen() error {
	err := c.client.ch.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("can't configure QoS: %w", err)
	}

	msgChannel, err := c.client.ch.Consume(c.client.queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("can't register channel: %w", err)
	}

	ctx := context.Background()

	for {
		msg := <-msgChannel
		logger.FromContext(ctx).Infoln("message received from queue")
		go func() {
			logger.FromContext(ctx).Infoln("message processing...")
			err := c.consumeFromQueue(ctx, &msg)
			if err != nil {
				nErr := msg.Nack(true, false)
				if nErr != nil {
					logger.FromContext(ctx).
						Errorln(fmt.Errorf("can't make negative acknowledgement: %w, (original error: %v)", nErr, err))
				}
				return
			}
			aErr := msg.Ack(false)
			if aErr != nil {
				logger.FromContext(ctx).Errorln(fmt.Errorf("can't make acknowledgement: %w", aErr))
			}
			logger.FromContext(ctx).Infoln("message processed successfully")
		}()
	}
}

// consumeFromQueue wraps the message processing and confirms its completion.
func (c *RabbitMQConsumer) consumeFromQueue(ctx context.Context, msg *amqp.Delivery) error {
	var data queueMessage
	err := json.NewDecoder(bytes.NewReader(msg.Body)).Decode(&data)
	if err != nil {
		return fmt.Errorf("can't decode queue message: %w", err)
	}

	err = c.process(ctx, data)
	if err != nil {
		uErr := c.requestsRepo.UpdateRequest(ctx, data.RequestID, repository.RequestStatusFailed, "")
		if uErr != nil {
			logger.FromContext(ctx).WithField("request_id", data.RequestID).
				Errorln(fmt.Errorf("can't update request: %w, (original error: %v)", err, err))
		}
		return fmt.Errorf("error processing a message from the queue: %w", err)
	}

	return nil
}

// process processes the current message from the queue.
func (c *RabbitMQConsumer) process(ctx context.Context, data queueMessage) error {
	sourceFile, err := c.storage.DownloadFile(data.FileID)
	if err != nil {
		return fmt.Errorf("storage error: %w", err)
	}
	logger.FromContext(ctx).WithField("file_id", data.FileID).
		Infoln("original file successfully downloaded from the S3 storage")

	targetFile, err := converter.Convert(sourceFile, data.TargetFormat, data.Ratio)
	if err != nil {
		return fmt.Errorf("converter error: %w", err)
	}
	logger.FromContext(ctx).WithField("file_id", data.FileID).
		Infoln("converter successfully processed the original file")

	err = c.requestsRepo.UpdateRequest(ctx, data.RequestID, repository.RequestStatusProcessing, "")
	if err != nil {
		return fmt.Errorf("can't update request with id %s: %w", data.RequestID, err)
	}
	logger.FromContext(ctx).WithField("request_id", data.RequestID).
		Infoln("request updated to the status \"processing\"")

	targetFileID, err := c.imagesRepo.InsertImage(ctx, data.Filename, data.TargetFormat)
	if err != nil {
		return fmt.Errorf("repository error: %w", err)
	}
	logger.FromContext(ctx).WithField("file_id", targetFileID).
		Infoln("converted file successfully saved in the database")

	err = c.storage.UploadFile(targetFile, targetFileID)
	if err != nil {
		return fmt.Errorf("storage error: %w", err)
	}
	logger.FromContext(ctx).WithField("file_id", targetFileID).
		Infoln("converted file successfully uploaded to the S3 storage")

	err = c.requestsRepo.UpdateRequest(ctx, data.RequestID, repository.RequestStatusDone, targetFileID)
	if err != nil {
		return fmt.Errorf("request updating error: %w", err)
	}
	logger.FromContext(ctx).WithField("request_id", data.RequestID).
		Infoln("request updated to the status \"done\"")

	return nil
}
