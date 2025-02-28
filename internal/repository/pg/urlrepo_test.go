package pg

import (
	"context"
	"database/sql"
	"net/http"
	"testing"

	"github.com/AlexNov03/UrlShortener/utils"
	"github.com/DATA-DOG/go-sqlmock"

	"github.com/stretchr/testify/assert"
)

func TestGetOriginalUrlOk(t *testing.T) {

	tests := []struct {
		Name          string
		ShortUrl      string
		Setup         func(m sqlmock.Sqlmock)
		ExpectOrigUrl string
		ExpectErr     error
		ExpectCode    int
	}{
		{
			Name:     "successful getting origUrl",
			ShortUrl: "Abc_efg",
			Setup: func(m sqlmock.Sqlmock) {
				rows := m.NewRows([]string{"original_url"}).AddRow(
					"http://ya.ru")
				m.ExpectQuery(`SELECT original_url FROM url WHERE short_url=\$1`).WithArgs(
					"Abc_efg").WillReturnRows(rows)
			},
			ExpectOrigUrl: "http://ya.ru",
			ExpectErr:     nil,
			ExpectCode:    0,
		},
		{
			Name:     "failed getting origUrl",
			ShortUrl: "Abc_efah",
			Setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT original_url FROM url WHERE short_url=\$1`).WithArgs(
					"Abc_efah").WillReturnError(sql.ErrNoRows)
			},
			ExpectOrigUrl: "",
			ExpectErr:     &utils.InternalError{Code: http.StatusNotFound, Message: "no originalUrl match this shortUrl"},
			ExpectCode:    http.StatusConflict,
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error while opening mock repo")
	}
	defer db.Close()

	urlRepo := NewUrlRepository(db)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			tt.Setup(mock)
			res, err := urlRepo.GetOriginalUrl(context.Background(), tt.ShortUrl)

			assert.Equal(t, res, tt.ExpectOrigUrl)
			assert.Equal(t, err, tt.ExpectErr)
		})
	}

}
