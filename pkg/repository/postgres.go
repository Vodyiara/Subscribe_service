package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	subscriptionsTable = "subscriptions"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgreDB(cfg Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	fmt.Println("DSN:", dsn)

	db, err := sqlx.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	fmt.Println("Running query on DB:", subscriptionsTable)
	if err != nil {
		return nil, err
	}

	return db, nil
}
