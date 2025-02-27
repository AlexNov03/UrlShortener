package delivery

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AlexNov03/UrlShortener/internal/models"
	"github.com/AlexNov03/UrlShortener/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type UrlUsecase interface {
	ShortenUrl(ctx context.Context, originalUrl string) (string, error)
	GetOriginalUrl(ctx context.Context, shortUrl string) (string, error)
}

type UrlDelivery struct {
	UC        UrlUsecase
	validator *validator.Validate
}

func NewUrlDelivery(uc UrlUsecase, validator *validator.Validate) *UrlDelivery {
	return &UrlDelivery{UC: uc, validator: validator}
}

func (ud *UrlDelivery) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	inputData := &models.OrigUrlData{}

	err := json.NewDecoder(r.Body).Decode(inputData)

	if err != nil {
		utils.ProcessBadRequestError(w, "delivery.ShortenUrl: incorrect input data")
		return
	}

	err = ud.validator.Struct(inputData)
	if err != nil {
		utils.ProcessBadRequestError(w, "delivery.ShortenUrl: incorrect fields in input data")
		return
	}

	ctx := r.Context()

	shortenedUrl, err := ud.UC.ShortenUrl(ctx, inputData.OriginalUrl)
	if err != nil {
		utils.ProcessError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&models.ShortUrlData{ShortUrl: shortenedUrl})

}

func (ud *UrlDelivery) GetOriginalUrl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl := vars["shortened_url"]

	ctx := r.Context()

	origUrl, err := ud.UC.GetOriginalUrl(ctx, shortUrl)
	if err != nil {
		utils.ProcessError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&models.OrigUrlData{OriginalUrl: origUrl})
}
