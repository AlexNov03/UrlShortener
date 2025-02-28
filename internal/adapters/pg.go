package adapters

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/AlexNov03/UrlShortener/internal/bootstrap"
	_ "github.com/lib/pq"
)

func GetDB(cfg *bootstrap.Config) (*sql.DB, error) {

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port,
		cfg.Database.DBName, cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
