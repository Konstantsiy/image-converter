// Package config works with application configurations.
package config

import "os"

// Config represents the application configurations.
type Config struct {
	Username        string
	Password        string
	DBName          string
	Host            string
	Port            string
	SSLMode         string
	PrivateKey      string
	BucketName      string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

// Load loads the necessary configurations from the .env file.
func (c *Config) Load() {
	c.Username = os.Getenv("DB_USERNAME")
	c.Password = os.Getenv("DB_PASSWORD")
	c.DBName = os.Getenv("DB_NAME")
	c.Host = os.Getenv("DB_HOST")
	c.Port = os.Getenv("DB_PORT")
	c.SSLMode = os.Getenv("DB_SSL_MODE")
	c.PrivateKey = os.Getenv("JWT_PRIVATE_KEY")
	c.BucketName = os.Getenv("AWS_BUCKET_NAME")
	c.Region = os.Getenv("AWS_REGION")
	c.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	c.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
}
