package usecase

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"math/rand"

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

	uc := NewUrlUsecase(mockRepo, rnd)
	ctx := context.Background()

	generatedShortUrl := uc.generateShortUrl()

	originalUrl := "http://example.ru"

	tests := []struct {
		Name           string
		SetUp          func()
		ExpectedString string
		ExpectedErr    error
	}{
		{
			Name: "Test for successful adding generated url",
			SetUp: func() {
				uc.rnd.Seed(64)
				mockRepo.EXPECT().GetOriginalUrl(ctx, generatedShortUrl).Return("", &utils.InternalError{
					Code: http.StatusNotFound, Message: "no originalUrl match this shortUrl"})

				mockRepo.EXPECT().AddOriginalUrl(ctx, &models.UrlData{OriginalUrl: originalUrl,
					ShortUrl: generatedShortUrl}).Return(nil)
			},
			ExpectedString: generatedShortUrl,
			ExpectedErr:    nil,
		},
		{
			Name: "Test for failed GetOriginalUrl request to db",
			SetUp: func() {
				uc.rnd.Seed(64)
				mockRepo.EXPECT().GetOriginalUrl(ctx, generatedShortUrl).Return("", fmt.Errorf("pg.UrlRepository.GetOriginalUrl:%w", context.DeadlineExceeded))
			},
			ExpectedString: "",
			ExpectedErr:    fmt.Errorf("pg.UrlRepository.GetOriginalUrl:%w", context.DeadlineExceeded),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			tt.SetUp()

			shortUrl, err := uc.ShortenUrl(ctx, originalUrl)

			assert.Equal(t, shortUrl, tt.ExpectedString)
			assert.Equal(t, err, tt.ExpectedErr)

		})
	}
}

func TestGetOriginalUrl(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUrlRepository(ctrl)

	rnd := rand.New(rand.NewSource(64))

	uc := NewUrlUsecase(mockRepo, rnd)
	ctx := context.Background()

	generatedShortUrl := uc.generateShortUrl()

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
				mockRepo.EXPECT().GetOriginalUrl(ctx, generatedShortUrl).Return(originalUrl, nil)
			},
			ExpectedString: originalUrl,
			ExpectedErr:    nil,
		},
		{
			Name: "Test for failed getting original url",
			SetUp: func() {
				mockRepo.EXPECT().GetOriginalUrl(ctx, generatedShortUrl).Return("",
					&utils.InternalError{Code: http.StatusNotFound, Message: "no originalUrl match this shortUrl"})
			},
			ExpectedString: "",
			ExpectedErr:    &utils.InternalError{Code: http.StatusNotFound, Message: "no originalUrl match this shortUrl"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			tt.SetUp()

			originalUrl, err := uc.GetOriginalUrl(ctx, generatedShortUrl)

			assert.Equal(t, originalUrl, tt.ExpectedString)
			assert.Equal(t, err, tt.ExpectedErr)

		})
	}
}
