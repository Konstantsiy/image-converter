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
	conf, err := config.Load()
	if err != nil {
		return fmt.Errorf("can't load configs: %w", err)
	}

	fmt.Printf("API configs: \n\tAppPort: %v\n\tDB: %+v\n\tJWT: %+v\n\tAWS: %+v\n\tRabbit: %+v\n",
		conf.AppPort, conf.DBConf, conf.JWTConf, conf.AWSConf, conf.RabbitMQConf)

	db, err := repository.NewPostgresDB(conf.DBConf)
	if err != nil {
		return fmt.Errorf("can't connect to postgres database: %v", err)
	}
	defer db.Close()
	logger.Info(context.Background(), fmt.Sprintf("database connected successfully (%s:%s)",
		conf.DBConf.Host, conf.DBConf.Port))

	repo := repository.NewRepository(db)

	st, err := storage.NewStorage(conf.AWSConf)
	if err != nil {
		return fmt.Errorf("can't create storage: %w", err)
	}
	logger.Info(context.Background(), "AWS S3 connected successfully")

	consumer, err := queue.NewConsumer(repo, st, conf.RabbitMQConf)
	if err != nil {
		return fmt.Errorf("can't create consumer: %w", err)
	}
	logger.Info(context.Background(), "RabbitMQ client (consumer) initialized successfully")

	err = consumer.Listen()
	if err != nil {
		return fmt.Errorf("can't start listening to the queue: %w", err)
	}
	logger.Info(context.Background(), "queue is listening")

	return nil
}
