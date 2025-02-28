package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AlexNov03/UrlShortener/internal/models"
	"github.com/AlexNov03/UrlShortener/utils"
)

type UrlRepository struct {
	DB *sql.DB
}

func NewUrlRepository(db *sql.DB) *UrlRepository {
	return &UrlRepository{DB: db}
}

func (ur *UrlRepository) AddOriginalUrl(ctx context.Context, data *models.UrlData) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := ur.DB.QueryRowContext(ctx, `SELECT 1 FROM url WHERE short_url=$1`, data.ShortUrl).Scan(new(int))
	if err == nil {
		return &utils.InternalError{Code: http.StatusConflict, Message: "this shortUrl already exists"}
	}

	if err != sql.ErrNoRows {
		return fmt.Errorf("pg.UrlRepository.AddOriginalUrl: %w", err)
	}

	_, err = ur.DB.ExecContext(ctx, `INSERT INTO url (short_url, original_url) VALUES ($1, $2)`, data.ShortUrl, data.OriginalUrl)
	if err != nil {
		return fmt.Errorf("pg.UrlRepository.AddOriginalUrl: %w", err)
	}

	return nil
}

func (ur *UrlRepository) GetOriginalUrl(ctx context.Context, shortUrl string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var originalUrl string
	err := ur.DB.QueryRowContext(ctx, `SELECT original_url FROM url WHERE short_url=$1`, shortUrl).Scan(&originalUrl)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", &utils.InternalError{Code: http.StatusNotFound, Message: "no originalUrl match this shortUrl"}
		}
		return "", fmt.Errorf("pg.UrlRepository.GetOriginalUrl: %w", err)
	}
	return originalUrl, nil
}
