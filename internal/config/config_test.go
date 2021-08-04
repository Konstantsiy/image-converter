package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	require := require.New(t)

	os.Setenv("DB_USERNAME", "postgres")
	os.Setenv("DB_PASSWORD", "qwerty123")
	os.Setenv("DB_NAME", "ita")
	os.Setenv("DB_HOST", "8080")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSL_MODE", "disable")

	os.Setenv("JWT_PUBLIC_KEY", "123456789")
	os.Setenv("JWT_PRIVATE_KEY", "1234567")

	os.Setenv("AWS_REGION", "eu-central-1")
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAUTWMM3GR4BUJVGPS")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")
	os.Setenv("AWS_BUCKET_NAME", "name1234")

	actual, err := Load()
	require.NoError(err)

	expected := Config{
		Username: "postgres",
		Password: "qwerty123",
		DBName:   "ita",
		Host:     "8080",
		Port:     "5432",
		SSLMode:  "disable",

		PublicKey:  "123456789",
		PrivateKey: "1234567",

		Region:          "eu-central-1",
		AccessKeyID:     "AKIAUTWMM3GR4BUJVGPS",
		SecretAccessKey: "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",
		BucketName:      "name1234",
	}
	require.Equal(expected, actual)
}
