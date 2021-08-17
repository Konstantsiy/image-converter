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
		return nil, err
	}

	return db, db.Ping()
}
