// Package app provides a function to start the application.
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Konstantsiy/image-converter/pkg/logger"

	"github.com/Konstantsiy/image-converter/internal/queue"

	"github.com/Konstantsiy/image-converter/internal/config"
	"github.com/Konstantsiy/image-converter/internal/repository"
	"github.com/Konstantsiy/image-converter/internal/server"
	"github.com/Konstantsiy/image-converter/internal/storage"
	"github.com/Konstantsiy/image-converter/pkg/jwt"
	"github.com/gorilla/mux"
)

func setEnv() {
	os.Setenv("APP_PORT", "8080")

	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "postgres")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSL_MODE", "disable")

	os.Setenv("JWT_PRIVATE_KEY_PATH", "./rsa_keys/private.pem")
	os.Setenv("JWT_PUBLIC_KEY_PATH", "./rsa_keys/public.pem")

	os.Setenv("AWS_REGION", "eu-central-1")
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAUTWMM3GR4BUJVGPS")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")
	os.Setenv("AWS_BUCKET_NAME", "name1234")

	os.Setenv("RABBITMQ_QUEUE_NAME", "new_queue")
	os.Setenv("RABBITMQ_AMQP_CONNECTION_URL", "amqp://guest:guest@localhost:5672/")
}

// Start starts the application server.
func Start() error {
	r := mux.NewRouter()
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

	tokenManager, err := jwt.NewTokenManager(conf.JWTConf)
	if err != nil {
		return fmt.Errorf("token manager error: %w", err)
	}
	logger.FromContext(context.Background()).Infoln("JWT-manager created successfully")

	st, err := storage.NewStorage(conf.AWSConf)
	if err != nil {
		return fmt.Errorf("can't create storage: %v", err)
	}
	logger.FromContext(context.Background()).Infoln("AWS S3 connected successfully")

	producer, err := queue.NewProducer(conf.RabbitMQConf)
	if err != nil {
		return fmt.Errorf("can't create producer: %w", err)
	}
	logger.FromContext(context.Background()).Infoln("RabbitMQ client (producer) initialized successfully")

	s := server.NewServer(repo, tokenManager, st, producer)
	s.RegisterRoutes(r)

	return http.ListenAndServe(":"+conf.AppPort, r)
}
