package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexNov03/UrlShortener/internal/delivery/mocks"
	"github.com/AlexNov03/UrlShortener/internal/models"
	"github.com/AlexNov03/UrlShortener/utils"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestShortenUrlOK(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedUc := mocks.NewMockUrlUsecase(ctrl)

	originalUrl := "http://ya.ru"
	shortUrl := "http://localhost:8080/Abc_def_qA"

	validator := validator.New(validator.WithRequiredStructEnabled())

	ud := NewUrlDelivery(mockedUc, validator)

	tests := []struct {
		Name                   string
		Setup                  func(context.Context)
		ReqBody                models.OrigUrlData
		ExpectedRespBody       models.ShortUrlData
		ExpectedRespStatusCode int
	}{
		{
			Name: "successful getting shorten url",
			Setup: func(ctx context.Context) {
				mockedUc.EXPECT().ShortenUrl(ctx, originalUrl).Return(shortUrl, nil)
			},
			ReqBody: models.OrigUrlData{
				OriginalUrl: originalUrl,
			},
			ExpectedRespBody: models.ShortUrlData{
				ShortUrl: shortUrl,
			},
			ExpectedRespStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			reqBodyJSON, err := json.Marshal(tt.ReqBody)
			assert.NoError(t, err)

			reqBody := bytes.NewReader(reqBodyJSON)

			r := httptest.NewRequest(http.MethodPost, "/shorten", reqBody)
			w := httptest.NewRecorder()

			tt.Setup(r.Context())

			ud.ShortenUrl(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.ExpectedRespStatusCode, resp.StatusCode)

			resultRespBody := models.ShortUrlData{}

			json.NewDecoder(resp.Body).Decode(&resultRespBody)

			assert.Equal(t, tt.ExpectedRespBody, resultRespBody)
		})
	}
}

func TestShortenUrlFail(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedUc := mocks.NewMockUrlUsecase(ctrl)

	originalUrl := "http://ya.ru"

	validator := validator.New(validator.WithRequiredStructEnabled())

	ud := NewUrlDelivery(mockedUc, validator)

	tests := []struct {
		Name                   string
		Setup                  func(context.Context)
		ReqBody                string
		ExpectedRespBody       utils.RestError
		ExpectedRespStatusCode int
	}{
		{
			Name:    "test for incorrect field names",
			Setup:   func(ctx context.Context) {},
			ReqBody: fmt.Sprintf(`{"original_ur":"%s"}`, originalUrl),
			ExpectedRespBody: utils.RestError{
				Error: "incorrect fields in input data",
			},
			ExpectedRespStatusCode: http.StatusBadRequest,
		},
		{
			Name: "test for incorrect original_url format",
			Setup: func(ctx context.Context) {
				mockedUc.EXPECT().ShortenUrl(ctx, "http:/localhost/ya.ru").Return("",
					utils.NewInternalError(http.StatusBadRequest, "original url does not fits the url format"),
				)
			},
			ReqBody: fmt.Sprintf(`{"original_url":"http:/localhost/ya.ru"}`),
			ExpectedRespBody: utils.RestError{
				Error: "original url does not fits the url format",
			},
			ExpectedRespStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			reqBody := bytes.NewReader([]byte(tt.ReqBody))

			r := httptest.NewRequest(http.MethodPost, "/shorten", reqBody)
			w := httptest.NewRecorder()

			tt.Setup(r.Context())

			ud.ShortenUrl(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.ExpectedRespStatusCode, resp.StatusCode)

			resultRespBody := utils.RestError{}

			json.NewDecoder(resp.Body).Decode(&resultRespBody)

			assert.Equal(t, tt.ExpectedRespBody, resultRespBody)
		})
	}
}

func TestGetOriginalUrlOK(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedUc := mocks.NewMockUrlUsecase(ctrl)

	originalUrl := "http://ya.ru"
	shortUrl := "/Abc_def_qA"
	shortUrlSuffix := "Abc_def_qA"

	validator := validator.New(validator.WithRequiredStructEnabled())

	ud := NewUrlDelivery(mockedUc, validator)

	router := mux.NewRouter()
	router.HandleFunc("/{shortened_url}", ud.GetOriginalUrl).Methods(http.MethodGet)

	tests := []struct {
		Name                   string
		Setup                  func(context.Context)
		ReqBody                models.ShortUrlData
		ExpectedRespBody       models.OrigUrlData
		ExpectedRespStatusCode int
	}{
		{
			Name: "successful getting original url",
			Setup: func(ctx context.Context) {
				mockedUc.EXPECT().GetOriginalUrl(gomock.Any(), shortUrlSuffix).Return(originalUrl, nil)
			},
			ReqBody: models.ShortUrlData{
				ShortUrl: shortUrl,
			},
			ExpectedRespBody: models.OrigUrlData{
				OriginalUrl: originalUrl,
			},
			ExpectedRespStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			r := httptest.NewRequest(http.MethodGet, shortUrl, nil)
			w := httptest.NewRecorder()

			tt.Setup(r.Context())

			router.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.ExpectedRespStatusCode, resp.StatusCode)

			resultRespBody := models.OrigUrlData{}

			json.NewDecoder(resp.Body).Decode(&resultRespBody)

			assert.Equal(t, tt.ExpectedRespBody, resultRespBody)
		})
	}
}

func TestGetOriginalUrlFail(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedUc := mocks.NewMockUrlUsecase(ctrl)

	shortUrl := "/Abc_def_qA"
	shortUrlSuffix := "Abc_def_qA"

	validator := validator.New(validator.WithRequiredStructEnabled())

	ud := NewUrlDelivery(mockedUc, validator)

	router := mux.NewRouter()
	router.HandleFunc("/{shortened_url}", ud.GetOriginalUrl).Methods(http.MethodGet)

	tests := []struct {
		Name                   string
		Setup                  func(context.Context)
		ReqBody                models.ShortUrlData
		ExpectedRespBody       utils.RestError
		ExpectedRespStatusCode int
	}{
		{
			Name: "error while getting original url",
			Setup: func(ctx context.Context) {
				mockedUc.EXPECT().GetOriginalUrl(gomock.Any(), shortUrlSuffix).Return("",
					utils.NewInternalError(http.StatusNotFound, "no originalUrl match this shortUrl"))
			},
			ReqBody: models.ShortUrlData{
				ShortUrl: shortUrl,
			},
			ExpectedRespBody: utils.RestError{
				Error: "no originalUrl match this shortUrl",
			},
			ExpectedRespStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			r := httptest.NewRequest(http.MethodGet, shortUrl, nil)
			w := httptest.NewRecorder()

			tt.Setup(r.Context())

			router.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.ExpectedRespStatusCode, resp.StatusCode)

			resultRespBody := utils.RestError{}

			json.NewDecoder(resp.Body).Decode(&resultRespBody)

			assert.Equal(t, tt.ExpectedRespBody, resultRespBody)
		})
	}
}
