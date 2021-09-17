// Package config works with application configurations.
package config

import (
	"github.com/kelseyhightower/envconfig"
)

type DBConfig struct {
	Username string `envconfig:"USERNAME"`
	Password string `envconfig:"PASSWORD"`
	DBName   string `envconfig:"NAME"`
	Host     string `envconfig:"HOST"`
	Port     string `envconfig:"PORT"`
	SSLMode  string `envconfig:"SSL_MODE"`
}

type AWSConfig struct {
	Region          string `envconfig:"REGION"`
	AccessKeyID     string `envconfig:"ACCESS_KEY_ID"`
	SecretAccessKey string `envconfig:"SECRET_ACCESS_KEY"`
	BucketName      string `envconfig:"BUCKET_NAME"`
}

type JWTConfig struct {
	PublicKeyPath  string `envconfig:"PUBLIC_KEY_PATH"`
	PrivateKeyPath string `envconfig:"PRIVATE_KEY_PATH"`
}

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
	err := envconfig.Process("envconfig", &c)
	return c, err
}
