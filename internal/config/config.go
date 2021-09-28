// Package config works with application configurations.
package config

import (
	"github.com/kelseyhightower/envconfig"
)

// DBConfig required to configure the database.
type DBConfig struct {
	User     string `envconfig:"USER"`
	Password string `envconfig:"PASSWORD"`
	DBName   string `envconfig:"NAME"`
	Host     string `envconfig:"HOST"`
	Port     string `envconfig:"PORT"`
	SSLMode  string `envconfig:"SSL_MODE"`
}

// AWSConfig required to configure the AWS S3 bucket.
type AWSConfig struct {
	Region          string `envconfig:"REGION"`
	AccessKeyID     string `envconfig:"ACCESS_KEY_ID"`
	SecretAccessKey string `envconfig:"SECRET_ACCESS_KEY"`
	BucketName      string `envconfig:"BUCKET_NAME"`
}

// JWTConfig required for configuring work with JWT.
type JWTConfig struct {
	PublicKeyPath  string `envconfig:"PUBLIC_KEY_PATH"`
	PrivateKeyPath string `envconfig:"PRIVATE_KEY_PATH"`
}

// RabbitMQConfig required to configure the RabbitMQ.
type RabbitMQConfig struct {
	QueueName         string `envconfig:"QUEUE_NAME"`
	AMQPConnectionURL string `envconfig:"AMQP_CONNECTION_URL"`
}

// Config represents the application configurations.
type Config struct {
	AppPort      string          `envconfig:"APP_PORT"`
	DBConf       *DBConfig       `envconfig:"DB"`
	JWTConf      *JWTConfig      `envconfig:"JWT"`
	AWSConf      *AWSConfig      `envconfig:"AWS"`
	RabbitMQConf *RabbitMQConfig `envconfig:"RABBITMQ"`
}

// Load loads the necessary configurations.
func Load() (Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	return c, err
}
