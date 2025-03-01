package pg

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/AlexNov03/UrlShortener/internal/models"
	"github.com/AlexNov03/UrlShortener/utils"
	"github.com/DATA-DOG/go-sqlmock"

	"github.com/stretchr/testify/assert"
)

func TestGetOriginalUrl(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	urlRepo := NewUrlRepository(db)

	tests := []struct {
		Name          string
		ShortUrl      string
		Setup         func(m sqlmock.Sqlmock)
		ExpectOrigUrl string
		ExpectErr     error
	}{
		{
			Name:     "successful getting origUrl",
			ShortUrl: "Abc_efg_ag",
			Setup: func(m sqlmock.Sqlmock) {
				rows := m.NewRows([]string{"original_url"}).AddRow(
					"http://ya.ru")
				m.ExpectQuery(`SELECT original_url FROM url WHERE short_url=\$1`).WithArgs(
					"Abc_efg_ag").WillReturnRows(rows)
			},
			ExpectOrigUrl: "http://ya.ru",
			ExpectErr:     nil,
		},
		{
			Name:     "failed getting origUrl",
			ShortUrl: "Abc_efah_a",
			Setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT original_url FROM url WHERE short_url=\$1`).WithArgs(
					"Abc_efah_a").WillReturnError(sql.ErrNoRows)
			},
			ExpectOrigUrl: "",
			ExpectErr:     &utils.InternalError{Code: http.StatusNotFound, Message: "no originalUrl match this shortUrl"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			tt.Setup(mock)
			res, err := urlRepo.GetOriginalUrl(context.Background(), tt.ShortUrl)

			assert.Equal(t, tt.ExpectOrigUrl, res)
			assert.Equal(t, tt.ExpectErr, err)
		})
	}

}

func TestAddOriginalUrl(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	urlRepo := NewUrlRepository(db)

	tests := []struct {
		Name      string
		UrlData   models.UrlData
		Setup     func(m sqlmock.Sqlmock, data models.UrlData)
		ExpectErr error
	}{
		{
			Name: "successful adding origUrl",
			UrlData: models.UrlData{
				ShortUrl:    "http://localhost:8080/Abc_def_gs",
				OriginalUrl: "http://ya.ru",
			},
			Setup: func(m sqlmock.Sqlmock, data models.UrlData) {

				m.ExpectQuery(`SELECT 1 FROM url WHERE short_url=\$1`).WithArgs(
					data.ShortUrl).WillReturnError(sql.ErrNoRows)

				m.ExpectExec(`INSERT INTO url \(short_url, original_url\) VALUES \(\$1, \$2\)`).WithArgs(
					data.ShortUrl, data.OriginalUrl).WillReturnResult(sqlmock.NewResult(1, 1))

			},
			ExpectErr: nil,
		},
		{
			Name: "error when origUrl already exists",
			UrlData: models.UrlData{
				ShortUrl:    "http://localhost:8080/Abc_def_gs",
				OriginalUrl: "http://ya.ru",
			},
			Setup: func(m sqlmock.Sqlmock, data models.UrlData) {

				rows := m.NewRows([]string{"1"}).AddRow(1)

				m.ExpectQuery(`SELECT 1 FROM url WHERE short_url=\$1`).WithArgs(
					data.ShortUrl).WillReturnRows(rows)

			},
			ExpectErr: &utils.InternalError{Code: http.StatusConflict, Message: "this shortUrl already exists"},
		},
		{
			Name: "internal db error test",
			UrlData: models.UrlData{
				ShortUrl:    "http://localhost:8080/Abc_def_gs",
				OriginalUrl: "http://ya.ru",
			},
			Setup: func(m sqlmock.Sqlmock, data models.UrlData) {

				m.ExpectQuery(`SELECT 1 FROM url WHERE short_url=\$1`).WithArgs(
					data.ShortUrl).WillReturnError(sql.ErrNoRows)

				m.ExpectExec(`INSERT INTO url \(short_url, original_url\) VALUES \(\$1, \$2\)`).WithArgs(
					data.ShortUrl, data.OriginalUrl).WillReturnError(fmt.Errorf("some bd error"))

			},
			ExpectErr: fmt.Errorf("pg.UrlRepository.AddOriginalUrl: %w", fmt.Errorf("some bd error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			tt.Setup(mock, tt.UrlData)

			err := urlRepo.AddOriginalUrl(context.Background(), &tt.UrlData)

			assert.Equal(t, tt.ExpectErr, err)
		})
	}

}
