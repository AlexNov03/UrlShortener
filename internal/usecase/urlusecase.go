package usecase

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/AlexNov03/UrlShortener/internal/models"
	"github.com/AlexNov03/UrlShortener/utils"
	"golang.org/x/exp/rand"
)

type UrlRepository interface {
	AddOriginalUrl(ctx context.Context, data *models.UrlData) error
	GetOriginalUrl(ctx context.Context, shortUrl string) (string, error)
}

type UrlUsecase struct {
	Repo UrlRepository
}

func NewUrlUsecase(repo UrlRepository) *UrlUsecase {
	return &UrlUsecase{Repo: repo}
}

const length = 10
const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func (uc *UrlUsecase) generateShortUrl() string {

	rand.Seed(uint64(time.Now().UnixNano()))
	res := strings.Builder{}
	for i := 0; i < length; i++ {
		res.WriteByte(charSet[rand.Intn(len(charSet))])
	}
	return res.String()

}

func (uc *UrlUsecase) ShortenUrl(ctx context.Context, originalUrl string) (string, error) {

	for {
		shortUrl := uc.generateShortUrl()
		_, err := uc.Repo.GetOriginalUrl(ctx, shortUrl)

		var interr *utils.InternalError

		if errors.As(err, &interr) && interr.Code == http.StatusNotFound {
			err = uc.Repo.AddOriginalUrl(ctx, &models.UrlData{OriginalUrl: originalUrl, ShortUrl: shortUrl})
			if err != nil {
				return "", err
			}
			return shortUrl, nil
		}

		if err != nil {
			return "", err
		}
	}
}

func (uc *UrlUsecase) GetOriginalUrl(ctx context.Context, shortUrl string) (string, error) {

	return uc.Repo.GetOriginalUrl(ctx, shortUrl)

}
