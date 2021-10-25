package app

import (
	"context"
	"fmt"

	"github.com/Konstantsiy/image-converter/pkg/logger"

	"github.com/Konstantsiy/image-converter/internal/config"
	"github.com/Konstantsiy/image-converter/internal/queue"
	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/internal/storage"
)

// StartListener starts the queue listener.
func StartListener() error {
	conf, err := config.Load()
	if err != nil {
		return fmt.Errorf("can't load configs: %w", err)
	}

	db, err := repository.NewPostgresDB(conf.DBConf)
	if err != nil {
		return fmt.Errorf("can't connect to postgres database: %v", err)
	}
	defer db.Close()
	logger.FromContext(context.Background()).
		WithField("host", conf.DBConf.Host).WithField("port", conf.DBConf.Port).
		Infoln("database connected successfully")

	st, err := storage.NewStorage(conf.AWSConf)
	if err != nil {
		return fmt.Errorf("can't create storage: %w", err)
	}
	logger.FromContext(context.Background()).Infoln("AWS S3 connected successfully")

	imageRepo, err := repository.NewImagesRepository(db)
	if err != nil {
		return fmt.Errorf("images repository creating error: %w", err)
	}

	requestsRepo, err := repository.NewRequestsRepository(db)
	if err != nil {
		return fmt.Errorf("requests repository creating error: %w", err)

	}

	consumer, err := queue.NewRabbitMQConsumer(requestsRepo, imageRepo, st, conf.RabbitMQConf)
	if err != nil {
		return fmt.Errorf("can't create consumer: %w", err)
	}
	logger.FromContext(context.Background()).Infoln("RabbitMQ client (consumer) initialized successfully")

	logger.FromContext(context.Background()).Infoln("queue is listening")
	err = consumer.Listen()
	if err != nil {
		return fmt.Errorf("can't start listening to the queue: %w", err)
	}

	return nil
}
