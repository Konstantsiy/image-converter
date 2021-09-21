package app

import (
	"fmt"

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

	db, err := repository.NewPostgresDB(conf.DBConf)
	if err != nil {
		return fmt.Errorf("can't connect to postgres database: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)

	st, err := storage.NewStorage(conf.AWSConf)
	if err != nil {
		return fmt.Errorf("can't create storage: %w", err)
	}

	consumer, err := queue.NewConsumer(repo, st, conf.RabbitMQConf)
	if err != nil {
		return fmt.Errorf("can't create consumer: %w", err)
	}

	err = consumer.Listen()
	if err != nil {
		return fmt.Errorf("can't start listening to the queue: %w", err)
	}

	return nil
}
