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

func StartListener() error {
	setEnv()
	conf, err := config.Load()
	if err != nil {
		return fmt.Errorf("can't load configs: %w", err)
	}

	db, err := repository.NewPostgresDB(conf.DBConf)
	if err != nil {
		return fmt.Errorf("can't connect to postgres database: %v", err)
	}
	defer db.Close()
	logger.FromContext(context.Background()).WithField("host", conf.DBConf.Host).
		Infoln("database connected successfully")

	repo := repository.NewRepository(db)

	st, err := storage.NewStorage(conf.AWSConf)
	if err != nil {
		return fmt.Errorf("can't create storage: %w", err)
	}
	logger.FromContext(context.Background()).Infoln("AWS S3 connected successfully")

	consumer, err := queue.NewConsumer(repo, st, conf.RabbitMQConf)
	if err != nil {
		return fmt.Errorf("can't create consumer: %w", err)
	}
	logger.FromContext(context.Background()).Infoln("RabbitMQ client (consumer) initialized successfully")

	err = consumer.Listen()
	if err != nil {
		return fmt.Errorf("can't start listening to the queue: %w", err)
	}
	logger.FromContext(context.Background()).Infoln("queue is listening")

	return nil
}
