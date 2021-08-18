package repository

import (
	"database/sql"
	"fmt"

	"github.com/Konstantsiy/image-converter/internal/config"
)

// NewPostgresDB opens new postgres connection.
func NewPostgresDB(c *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.DBName, c.Password, c.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("error when opening a database connection: %v", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error when trying to ping the database: %v", err)
	}

	return db, nil
}
