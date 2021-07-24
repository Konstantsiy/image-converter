package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLoad tests Load function.
func TestLoad(t *testing.T) {
	require := require.New(t)

	os.Setenv("DB_USERNAME", "postgres")
	os.Setenv("DB_PASSWORD", "qwerty123")
	os.Setenv("DB_NAME", "ita")
	os.Setenv("DB_HOST", "8080")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSL_MODE", "disable")
	os.Setenv("JWT_PUBLIC_KEY", "123456789")

	actual, err := Load()
	require.NoError(err)

	expected := Config{
		Username:  "postgres",
		Password:  "qwerty123",
		DBName:    "ita",
		Host:      "8080",
		Port:      "5432",
		SSLMode:   "disable",
		PublicKey: "123456789",
	}
	require.Equal(expected, actual)
}
