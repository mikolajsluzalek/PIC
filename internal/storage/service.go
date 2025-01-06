package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/pkg/errors"
)

type Service struct {
	DB *sql.DB
}

func New() (Service, error) {
	svc := Service{}

	fmt.Println("Initializing database connection...")

	cfg, err := readConfig()
	if err != nil {
		return svc, errors.Wrap(err, "failed to read config")
	}

	db, err := sql.Open("sqlserver", cfg.DatabaseURL)
	if err != nil {
		return svc, errors.Wrap(err, "failed to open database")
	}

	svc.DB = db

	// Health check of the database connection
	err = db.Ping()
	if err != nil {
		return svc, errors.Wrap(err, "failed to ping database")
	}

	fmt.Println("Database connection established successfully!")

	return svc, nil
}
