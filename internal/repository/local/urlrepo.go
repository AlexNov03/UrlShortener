package local

import (
	"context"
	"net/http"
	"sync"

	"github.com/AlexNov03/UrlShortener/internal/models"
	"github.com/AlexNov03/UrlShortener/utils"
)

type UrlRepository struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewUrlRepository() *UrlRepository {
	return &UrlRepository{mu: sync.RWMutex{}, store: make(map[string]string)}
}

func (ur *UrlRepository) AddOriginalUrl(ctx context.Context, data *models.UrlData) error {

	origUrl := data.OriginalUrl
	shortUrl := data.ShortUrl

	ur.mu.Lock()
	defer ur.mu.Unlock()

	if _, ok := ur.store[shortUrl]; ok {
		return &utils.InternalError{Code: http.StatusConflict, Message: "this shortUrl already exists"}
	}
	ur.store[shortUrl] = origUrl
	return nil
}

func (ur *UrlRepository) GetOriginalUrl(ctx context.Context, shortUrl string) (string, error) {

	ur.mu.RLock()
	defer ur.mu.RUnlock()

	val, ok := ur.store[shortUrl]
	if !ok {
		return "", &utils.InternalError{Code: http.StatusNotFound, Message: "no originalUrl match this shortUrl"}
	}
	return val, nil
}
