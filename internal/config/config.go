// Package config works with application configurations.
package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents the application configurations.
type Config struct {
	Username   string `envconfig:"DB_USERNAME"`
	Password   string `envconfig:"DB_PASSWORD"`
	DBName     string `envconfig:"DB_NAME"`
	Host       string `envconfig:"DB_HOST"`
	Port       string `envconfig:"DB_PORT"`
	SSLMode    string `envconfig:"DB_SSL_MODE"`
	PublicKey  string `envconfig:"JWT_PUBLIC_KEY"`
	PrivateKey string `envconfig:"JWT_PRIVATE_KEY"`
}

// Load loads the necessary configurations.
func Load() (Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	return c, err
}
