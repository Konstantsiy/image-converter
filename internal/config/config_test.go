package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	require := require.New(t)

	os.Setenv("APP_PORT", "8080")

	os.Setenv("DB_USERNAME", "postgres")
	os.Setenv("DB_PASSWORD", "qwerty123")
	os.Setenv("DB_NAME", "ita")
	os.Setenv("DB_HOST", "8080")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSL_MODE", "disable")

	os.Setenv("JWT_PUBLIC_KEY_PATH", "123456789")
	os.Setenv("JWT_PRIVATE_KEY_PATH", "1234567")

	os.Setenv("AWS_REGION", "eu-central-1")
	os.Setenv("AWS_ACCESS_KEY_ID", "SGFHSGDHFSGF")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "SDFSDFDSFSF84378FDSFSDFSDFD")
	os.Setenv("AWS_BUCKET_NAME", "name1234")

	os.Setenv("RABBITMQ_QUEUE_NAME", "converter_queue")
	os.Setenv("RABBITMQ_CONNECTION_URL", "amqp://guest:guest@localhost:5672/")

	actual, err := Load()
	require.NoError(err)

	expected := Config{
		AppPort: "8080",

		Username: "postgres",
		Password: "qwerty123",
		DBName:   "ita",
		Host:     "8080",
		Port:     "5432",
		SSLMode:  "disable",

		PublicKeyPath:  "123456789",
		PrivateKeyPath: "1234567",

		Region:          "eu-central-1",
		AccessKeyID:     "SGFHSGDHFSGF",
		SecretAccessKey: "SDFSDFDSFSF84378FDSFSDFSDFD",
		BucketName:      "name1234",

		QueueName:         "converter_queue",
		AMQPConnectionURL: "amqp://guest:guest@localhost:5672/",
	}
	require.Equal(expected, actual)
}
