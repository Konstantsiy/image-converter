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

	actual, err := Load()
	require.NoError(err)

	expected := Config{
		Username:   "postgres",
		Password:   "qwerty123",
		DBName:     "ita",
		Host:       "8080",
		Port:       "5432",
		SSLMode:    "disable",
		PublicKey:  "123456789",
		PrivateKey: "1234567",
	}
	require.Equal(expected, actual)
}
