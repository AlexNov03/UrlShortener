package usecase

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"math/rand"

	"github.com/AlexNov03/UrlShortener/internal/bootstrap"
	"github.com/AlexNov03/UrlShortener/internal/models"
	"github.com/AlexNov03/UrlShortener/internal/usecase/mocks"
	"github.com/AlexNov03/UrlShortener/utils"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestShortenUrl(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUrlRepository(ctrl)

	rnd := rand.New(rand.NewSource(64))

	cfg := &bootstrap.Config{}
	cfg.Server.Protocol, cfg.Server.Host, cfg.Server.Port = "http", "localhost", 8080

	uc := NewUrlUsecase(mockRepo, rnd, cfg)
	ctx := context.Background()

	suffix := uc.generateShortUrl()
	generatedShortUrl := fmt.Sprintf("%s://%s:%d/%s", uc.cfg.Server.Protocol, uc.cfg.Server.Host, uc.cfg.Server.Port, suffix)

	tests := []struct {
		Name           string
		OriginalUrl    string
		SetUp          func(string)
		ExpectedString string
		ExpectedErr    error
	}{
		{
			Name:        "Test for successful returning generated url",
			OriginalUrl: "http://ozon.ru",
			SetUp: func(originalUrl string) {

				mockRepo.EXPECT().GetShortUrlByLong(ctx, originalUrl).Return("",
					&utils.InternalError{Code: http.StatusNotFound, Message: "no shortUrl match this originalUrl"})

				mockRepo.EXPECT().GetOriginalUrl(ctx, suffix).Return("", &utils.InternalError{
					Code: http.StatusNotFound, Message: "no originalUrl match this shortUrl"})

				mockRepo.EXPECT().AddOriginalUrl(ctx, &models.UrlData{OriginalUrl: originalUrl,
					ShortUrl: suffix}).Return(nil)
			},
			ExpectedString: generatedShortUrl,
			ExpectedErr:    nil,
		},
		{
			Name:        "Test for failed GetOriginalUrl request to db",
			OriginalUrl: "http://ozon.ru",
			SetUp: func(originalUrl string) {
				mockRepo.EXPECT().GetShortUrlByLong(ctx, originalUrl).Return("",
					&utils.InternalError{Code: http.StatusNotFound, Message: "no shortUrl match this originalUrl"})

				mockRepo.EXPECT().GetOriginalUrl(ctx, suffix).Return("", fmt.Errorf("pg.UrlRepository.GetOriginalUrl:%w", context.DeadlineExceeded))
			},
			ExpectedString: "",
			ExpectedErr:    fmt.Errorf("pg.UrlRepository.GetOriginalUrl:%w", context.DeadlineExceeded),
		},
		{
			Name:        "Test for bad url format",
			OriginalUrl: "httpozon.ru",
			SetUp: func(originalUrl string) {
				mockRepo.EXPECT().GetShortUrlByLong(ctx, originalUrl).Return("",
					&utils.InternalError{Code: http.StatusNotFound, Message: "no shortUrl match this originalUrl"})
			},
			ExpectedString: "",
			ExpectedErr:    &utils.InternalError{Code: http.StatusBadRequest, Message: "original url does not fits the url format"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			uc.rnd.Seed(64)
			tt.SetUp(tt.OriginalUrl)

			shortUrl, err := uc.ShortenUrl(ctx, tt.OriginalUrl)

			assert.Equal(t, tt.ExpectedString, shortUrl)
			assert.Equal(t, tt.ExpectedErr, err)

		})
	}
}

func TestGetOriginalUrl(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUrlRepository(ctrl)

	rnd := rand.New(rand.NewSource(64))

	cfg := &bootstrap.Config{}

	uc := NewUrlUsecase(mockRepo, rnd, cfg)
	ctx := context.Background()

	suffix := uc.generateShortUrl()

	originalUrl := "http://example.ru"

	tests := []struct {
		Name           string
		SetUp          func()
		ExpectedString string
		ExpectedErr    error
	}{
		{
			Name: "Test for successful getting original url",
			SetUp: func() {
				mockRepo.EXPECT().GetOriginalUrl(ctx, suffix).Return(originalUrl, nil)
			},
			ExpectedString: originalUrl,
			ExpectedErr:    nil,
		},
		{
			Name: "Test for failed getting original url",
			SetUp: func() {
				mockRepo.EXPECT().GetOriginalUrl(ctx, suffix).Return("",
					&utils.InternalError{Code: http.StatusNotFound, Message: "no originalUrl match this shortUrl"})
			},
			ExpectedString: "",
			ExpectedErr:    &utils.InternalError{Code: http.StatusNotFound, Message: "no originalUrl match this shortUrl"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			tt.SetUp()

			originalUrl, err := uc.GetOriginalUrl(ctx, suffix)

			assert.Equal(t, tt.ExpectedString, originalUrl)
			assert.Equal(t, tt.ExpectedErr, err)

		})
	}
}
