package local

import (
	"context"
	"net/http"
	"sync"

	"github.com/AlexNov03/UrlShortener/internal/models"
	"github.com/AlexNov03/UrlShortener/utils"
)

type UrlRepository struct {
	mu            sync.RWMutex
	storeShortUrl map[string]string
	storeLongUrl  map[string]string
}

func NewUrlRepository() *UrlRepository {
	return &UrlRepository{mu: sync.RWMutex{}, storeShortUrl: make(map[string]string), storeLongUrl: make(map[string]string)}
}

func (ur *UrlRepository) AddOriginalUrl(ctx context.Context, data *models.UrlData) error {

	origUrl := data.OriginalUrl
	shortUrl := data.ShortUrl

	ur.mu.Lock()
	defer ur.mu.Unlock()

	if _, ok := ur.storeShortUrl[shortUrl]; ok {
		return utils.NewInternalError(http.StatusConflict, "this shortUrl already exists")
	}
	if _, ok := ur.storeLongUrl[origUrl]; ok {
		return utils.NewInternalError(http.StatusConflict, "this longUrl already exists")
	}

	ur.storeShortUrl[shortUrl] = origUrl
	ur.storeLongUrl[origUrl] = shortUrl

	return nil
}

func (ur *UrlRepository) GetOriginalUrl(ctx context.Context, shortUrl string) (string, error) {

	ur.mu.RLock()
	defer ur.mu.RUnlock()

	val, ok := ur.storeShortUrl[shortUrl]
	if !ok {
		return "", &utils.InternalError{Code: http.StatusNotFound, Message: "no originalUrl match this shortUrl"}
	}
	return val, nil
}

func (ur *UrlRepository) GetShortUrlByLong(ctx context.Context, longUrl string) (string, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	val, ok := ur.storeLongUrl[longUrl]
	if !ok {
		return "", &utils.InternalError{Code: http.StatusNotFound, Message: "no shortUrl match this originalUrl"}
	}

	return val, nil
}
