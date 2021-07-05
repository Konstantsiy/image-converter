package config

import "os"

type Config struct {
	Username   string
	Password   string
	DBName     string
	Host       string
	Port       string
	SSLMode    string
	PrivateKey string
}

func (c *Config) Load() {
	c.Username = os.Getenv("DB_USERNAME")
	c.Password = os.Getenv("DB_PASSWORD")
	c.DBName = os.Getenv("DB_NAME")
	c.Host = os.Getenv("DB_HOST")
	c.Port = os.Getenv("DB_PORT")
	c.SSLMode = os.Getenv("DB_SSL_MODE")
	c.PrivateKey = os.Getenv("JWT_PRIVATE_KEY")
}
