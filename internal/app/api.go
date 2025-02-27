package app

import (
	"github.com/AlexNov03/UrlShortener/internal/bootstrap"
	"github.com/AlexNov03/UrlShortener/internal/delivery"
	localrepo "github.com/AlexNov03/UrlShortener/internal/repository/local"
	"github.com/AlexNov03/UrlShortener/internal/server"
	"github.com/AlexNov03/UrlShortener/internal/usecase"
	"github.com/go-playground/validator/v10"
)

type ApiEntryPoint struct {
	cfg    *bootstrap.Config
	server *server.Server
}

func NewApiEntryPoint() *ApiEntryPoint {
	return &ApiEntryPoint{}
}

func (ae *ApiEntryPoint) Init() error {
	config, err := bootstrap.ReadConfig()
	if err != nil {
		return err
	}

	ae.cfg = config

	validator := validator.New(validator.WithRequiredStructEnabled())

	repo := localrepo.NewUrlRepository()
	uc := usecase.NewUrlUsecase(repo)
	deliv := delivery.NewUrlDelivery(uc, validator)

	ae.server = server.NewServer(ae.cfg, deliv)
	ae.server.Init()

	return nil
}

func (ae *ApiEntryPoint) Run() error {
	return ae.server.Run()

}

func (ae *ApiEntryPoint) Stop() error {
	return ae.server.Stop()
}
